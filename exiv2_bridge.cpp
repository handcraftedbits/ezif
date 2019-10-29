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

void readImageMetadata (const char *filename, exiv2Error *err, valueHolder *vh, readHandlers *handlers,
     void *rhPointer)
{
     try
     {
          std::ostringstream buffer;
          auto image = Exiv2::ImageFactory::open(filename);

          image->readMetadata();

          // Read Exif metadata.

          auto exifData = image->exifData();

          exifData.sortByKey();

          for (auto &exifDatum : exifData)
          {
               long count = getAdjustedCount(exifDatum.typeId(), exifDatum.count());
               std::string interpretedValue;

               buffer.clear();
               buffer.str("");

               // The interpreted value is obtained via operator<<.

               buffer << exifDatum;

               interpretedValue = buffer.str();

               // Notify that new Exif data has been encountered.

               handlers->dosc(rhPointer, exifDatum.familyName(), exifDatum.groupName().c_str(),
                    exifDatum.tagName().c_str(), (int) exifDatum.typeId(), exifDatum.tagLabel().c_str(),
                    interpretedValue.c_str(), count);

               // Notify that a value component has been encountered.

               for (int i = 0; i < count; ++i)
               {
                    populateValueHolder(vh, exifDatum.value(), i);

                    handlers->vc(rhPointer, exifDatum.familyName(), vh);
               }
          }
     }

     catch (Exiv2::Error &e)
     {
          err->code = e.code();
          err->message = strdup(e.what());
     }
}
