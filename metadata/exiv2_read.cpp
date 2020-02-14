#include <cstring>
#include <sstream>
#include <vector>

#include <exiv2/basicio.hpp>
#include <exiv2/exiv2.hpp>

#include "exiv2.h"

long getAdjustedCount (Exiv2::TypeId typeId, long count)
{
     switch (typeId)
     {
          // Exiv2 treats string values as an array of characters, but we're interested in the value as a single string,
          // string, so we'll adjust the count accordingly.

          case Exiv2::TypeId::asciiString:
          case Exiv2::TypeId::comment:
          case Exiv2::TypeId::string:
          case Exiv2::TypeId::xmpText:
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

void notifyValueCreated (valueHolder *vh, const Exiv2::Value &value, int index, readHandler *handler, void *rhPointer)
{
     vh->strValue = NULL;

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
               vh->strValue = strdup(value.toString(index).c_str());

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

          case Exiv2::TypeId::langAlt:
          {
               auto langAltMap = static_cast<const Exiv2::LangAltValue&>(value).value_;

               // XMPLangAlt is a little different from normal values since it's a map.  We have to iterate over it and
               // notify that a new value has been created for each key/value pair.

               for (auto langAlt : langAltMap)
               {
                    vh->langValue = langAlt.first.c_str();
                    vh->strValue = strdup(langAlt.second.c_str());

                    handler->vc(rhPointer, vh);

                    free(vh->strValue);
               }

               return;
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
     }

     handler->vc(rhPointer, vh);

     if (vh->strValue)
     {
          free(vh->strValue);
     }
}

void handleMetadatum (const Exiv2::Metadatum& metadatum, std::ostringstream& buffer, int repeatable, valueHolder *vh,
     readHandler *handler, void *rhPointer)
{
     long count = getAdjustedCount(metadatum.typeId(), metadatum.count());
     std::string interpretedValue;

     buffer.clear();
     buffer.str("");

     // The interpreted value can only(?) be obtained via operator<<.

     buffer << metadatum;

     interpretedValue = buffer.str();

     // Notify that new metadata has been encountered.

     handler->posc(rhPointer, metadatum.familyName(), metadatum.groupName().c_str(), metadatum.tagName().c_str(),
          (int) metadatum.typeId(), metadatum.tagLabel().c_str(), interpretedValue.c_str(), count, repeatable);

     for (int i = 0; i < count; ++i)
     {
          notifyValueCreated(vh, metadatum.value(), i, handler, rhPointer);
     }

     // Notify that we've finished processing the metadata.

     handler->poec(rhPointer, metadatum.familyName());
}

void readMetadata (Exiv2::Image::AutoPtr image, exiv2Error *err, valueHolder *vh, readHandler *handler, void *rhPointer)
{
     try
     {
          std::ostringstream buffer;

          image->readMetadata();

          for (auto &exifDatum : image->exifData())
          {
               handleMetadatum(exifDatum, buffer, 0, vh, handler, rhPointer);
          }

          for (auto &iptcDatum : image->iptcData())
          {
               handleMetadatum(iptcDatum, buffer, Exiv2::IptcDataSets::dataSetRepeatable(iptcDatum.tag(),
                    iptcDatum.record()) ? 1 : 0, vh, handler, rhPointer);
          }

          for (auto &xmpDatum : image->xmpData())
          {
               handleMetadatum(xmpDatum, buffer, 0, vh, handler, rhPointer);
          }
     }

     catch (Exiv2::Error &e)
     {
          err->code = e.code();
          err->message = strdup(e.what());
     }
}

void readCollectionFromFile (const char *filename, exiv2Error *err, valueHolder *vh, readHandler *handler,
     void *rhPointer)
{
     Exiv2::BasicIo::AutoPtr ptr(new Exiv2::FileIo(std::string(filename)));

     readMetadata(Exiv2::ImageFactory::open(ptr), err, vh, handler, rhPointer);
}

void readCollectionFromURL (const char *url, exiv2Error *err, valueHolder *vh, readHandler *handler, void *rhPointer)
{
     Exiv2::BasicIo::AutoPtr ptr(new Exiv2::HttpIo(std::string(url)));

     readMetadata(Exiv2::ImageFactory::open(ptr), err, vh, handler, rhPointer);
}
