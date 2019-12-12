#include <cstring>
#include <sstream>
#include <vector>

#include <exiv2/exiv2.hpp>

#include "exiv2_bridge.h"

long getAdjustedCount (Exiv2::TypeId typeId, long count)
{
     switch (typeId)
     {
          // Exiv2 treats Exif string values as an array of characters, but we're interested in the value as a single
          // string, so we'll adjust the count accordingly.

          case Exiv2::TypeId::asciiString:
          case Exiv2::TypeId::comment:
          case Exiv2::TypeId::string:
          {
               return 1;
          }

          // IPTC date and time values also seem to have a count that's not what we expect (1).

          case Exiv2::TypeId::date:
          case Exiv2::TypeId::time:
          {
               return 1;
          }
     }

     return count;
}

void populateValueHolder (valueHolder *vh, const Exiv2::Value &value, int index)
{
     switch (value.typeId())
     {
          case Exiv2::TypeId::asciiString:
          case Exiv2::TypeId::comment:
          case Exiv2::TypeId::string:
          case Exiv2::TypeId::xmpAlt:
          case Exiv2::TypeId::xmpBag:
          case Exiv2::TypeId::xmpSeq:
          case Exiv2::TypeId::xmpText:
          {
               vh->strValue = static_cast<const Exiv2::StringValueBase&>(value).value_.c_str();

               break;
          }

          case Exiv2::TypeId::date:
          {
               auto date = static_cast<const Exiv2::DateValue&>(value).getDate();

               vh->dayValue = date.day;
               vh->monthValue = date.month;
               vh->yearValue = date.year;

               break;
          }

          case Exiv2::TypeId::time:
          {
               auto time = static_cast<const Exiv2::TimeValue&>(value).getTime();

               vh->hourValue = time.hour;
               vh->minuteValue = time.minute;
               vh->secondValue = time.second;
               vh->timezoneHourOffset = time.tzHour;
               vh->timezoneMinuteOffset = time.tzMinute;

               break;
          }

          case Exiv2::TypeId::signedByte:
          case Exiv2::TypeId::signedLong:
          case Exiv2::TypeId::signedShort:
          case Exiv2::TypeId::undefined:
          case Exiv2::TypeId::unsignedByte:
          case Exiv2::TypeId::unsignedLong:
          case Exiv2::TypeId::unsignedShort:
          {
               vh->longValue = value.toLong(index);

               break;
          }

          case Exiv2::TypeId::signedRational:
          case Exiv2::TypeId::unsignedRational:
          {
               auto rationalValue = value.toRational(index);

               vh->rationalValueN = rationalValue.first;
               vh->rationalValueD = rationalValue.second;

               break;
          }

          case Exiv2::TypeId::tiffDouble:
          {
               // The Exiv2::Value class doesn't seem to have a function for getting a double value, just a float, so we
               // have to read it in the hard way.

               auto doubleValue = static_cast<const Exiv2::ValueType<double>&>(value);
               auto doubleValues = static_cast<std::vector<double>>(doubleValue.value_);

               vh->doubleValue = doubleValues[index];

               break;
          }

          case Exiv2::TypeId::tiffFloat:
          {
               vh->doubleValue = value.toFloat(index);

               break;
          }
     }
}

void readMetadatum (const Exiv2::Metadatum& metadatum, std::ostringstream& buffer, int repeatable, valueHolder *vh,
     readHandlers *handlers, void *rhPointer)
{
     long count = getAdjustedCount(metadatum.typeId(), metadatum.count());
     std::string interpretedValue;

     buffer.clear();
     buffer.str("");

     // The interpreted value can only(?) be obtained via operator<<.

     buffer << metadatum;

     interpretedValue = buffer.str();

     // Notify that new IPTC data has been encountered.

     handlers->dosc(rhPointer, metadatum.familyName(), metadatum.groupName().c_str(), metadatum.tagName().c_str(),
          (int) metadatum.typeId(), metadatum.tagLabel().c_str(), interpretedValue.c_str(), count, repeatable);

     // Notify that a value component has been encountered.

     for (int i = 0; i < count; ++i)
     {
          populateValueHolder(vh, metadatum.value(), i);

          handlers->vc(rhPointer, vh);
     }

     // Notify that we've finished processing the Exif data.

     handlers->doec(rhPointer, metadatum.familyName());
}

void readImageMetadata (const char *filename, exiv2Error *err, valueHolder *vh, readHandlers *handlers, void *rhPointer)
{
     try
     {
          std::ostringstream buffer;
          auto image = Exiv2::ImageFactory::open(filename);

          image->readMetadata();

          for (auto &exifDatum : image->exifData())
          {
               readMetadatum(exifDatum, buffer, 0, vh, handlers, rhPointer);
          }

          for (auto &iptcDatum : image->iptcData())
          {
               readMetadatum(iptcDatum, buffer, Exiv2::IptcDataSets::dataSetRepeatable(iptcDatum.tag(),
                    iptcDatum.record()) ? 1 : 0, vh, handlers, rhPointer);
          }
     }

     catch (Exiv2::Error &e)
     {
          err->code = e.code();
          err->message = strdup(e.what());
     }
}

