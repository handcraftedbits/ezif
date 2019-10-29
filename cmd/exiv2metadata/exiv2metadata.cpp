#include <cstdlib>
#include <iostream>
#include <set>

#include <exiv2/exiv2.hpp>
#include <jansson.h>

#include "exiv2metadata.h"

json_t *dumpExifTagToJSON (const Exiv2::TagInfo *tagInfo)
{
     return json_pack("{s:s, s:s, s:i, s:i}", "label", tagInfo->title_, "description", tagInfo->desc_, "typeId",
          tagInfo->typeId_, "count", tagInfo->count_);
}

json_t *dumpExifGroupToJSON (const char *groupName)
{
     const Exiv2::TagInfo *tags = Exiv2::ExifTags::tagList(groupName);
     json_t *root = json_object();

     // Not documented.  Exiv2 tag arrays use 0xFFFF as the tag ID of the last element.

     while (tags->tag_ != 0xFFFF)
     {
          json_object_set(root, tags->name_, dumpExifTagToJSON(tags));

          ++tags;
     }

     return root;
}

json_t *dumpExifGroupsToJSON (void)
{
     const Exiv2::GroupInfo *groups = Exiv2::ExifTags::groupList();
     json_t *root = json_object();
     std::set<Exiv2::TagListFct> seenGroups;

     // Not documented.  The Exiv2 Exif group array indicates the end of the array with a NULL tagList_ function.

     while (groups->tagList_ != NULL)
     {
          // Some groups are aliases, so we don't want to bother adding those tags again.  The only way to reliably
          // determine this is to see if the group has the same tag list function as another one.

          if (seenGroups.find(groups->tagList_) == seenGroups.end())
          {
               json_object_set(root, groups->groupName_, dumpExifGroupToJSON(groups->groupName_));

               seenGroups.insert(groups->tagList_);
          }

          ++groups;
     }

     return root;
}

json_t *dumpIPTCDataSetRecordToJSON (const Exiv2::DataSet *dataSet)
{
     return json_pack("{s:s, s:s, s:i, s:b, s:i, s:i}", "label", dataSet->title_, "description", dataSet->desc_,
          "typeId", dataSet->type_, "repeatable", dataSet->repeatable_, "minBytes", dataSet->minbytes_, "maxBytes",
          dataSet->maxbytes_);
}

json_t *dumpIPTCDataSetToJSON (const Exiv2::DataSet *dataSet)
{
     json_t *root = json_object();

     // Not documented.  Exiv2 IPTC data set arrays use 0xFFFF as the record ID of the last element.

     while (dataSet->number_ != 0xFFFF)
     {
          json_object_set(root, dataSet->name_, dumpIPTCDataSetRecordToJSON(dataSet));

          ++dataSet;
     }

     return root;
}

json_t *dumpIPTCDataSetsToJSON (void)
{
     json_t *root = json_object();

     json_object_set(root, Exiv2::IptcDataSets::recordName(Exiv2::IptcDataSets::envelope).c_str(),
          dumpIPTCDataSetToJSON(Exiv2::IptcDataSets::envelopeRecordList()));
     json_object_set(root, Exiv2::IptcDataSets::recordName(Exiv2::IptcDataSets::application2).c_str(),
          dumpIPTCDataSetToJSON(Exiv2::IptcDataSets::application2RecordList()));

     return root;
}

json_t *dumpXMPPropertyToJSON (const Exiv2::XmpPropertyInfo *property)
{
    return json_pack("{s:s, s:s, s:i}", "label", property->title_, "description", property->desc_, "typeId",
          property->typeId_);
}

json_t *dumpXMPNamespaceToJSON (const Exiv2::XmpNsInfo *ns)
{
     const Exiv2::XmpPropertyInfo *properties = ns->xmpPropertyInfo_;
     json_t *root;

     // Some entries in Exiv2's mappings are null apparently...

     if (properties == NULL)
     {
          return NULL;
     }

     root = json_object();

     while (properties->typeId_ != Exiv2::invalidTypeId)
     {
          json_object_set(root, properties->name_, dumpXMPPropertyToJSON(properties));

          ++properties;
     }

     return root;
}

json_t *dumpXMPNamespacesToJSON (void)
{
     Exiv2::Dictionary dict;
     json_t *root = json_object();

     Exiv2::XmpProperties::registeredNamespaces(dict);

     for (auto nsMapping : dict)
     {
          try
          {
               json_t *namespaceJSON = dumpXMPNamespaceToJSON(Exiv2::XmpProperties::nsInfo(nsMapping.first));

               if (namespaceJSON != NULL)
               {
                    json_object_set(root, nsMapping.first.c_str(), namespaceJSON);
               }
          }

          catch (Exiv2::Error &e)
          {
               // Ignore.  The list of namespaces seems to include ones from the XMP SDK for which there are no
               // associated XmpNsInfo objects, and other oddities.
          }
     }

     return root;
}

void dumpMetadataToJSON (void)
{
     char *output;
     json_t *root = json_object();

     json_object_set(root, "Exif", dumpExifGroupsToJSON());
     json_object_set(root, "IPTC", dumpIPTCDataSetsToJSON());
     json_object_set(root, "XMP", dumpXMPNamespacesToJSON());

     output = json_dumps(root, JSON_INDENT(2) | JSON_SORT_KEYS);

     std::cout << output << std::endl;

     free(output);

     json_decref(root);
}
