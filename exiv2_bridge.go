package ezif

/*
#include "exiv2_bridge.h"

extern void onDatumStartGo(void*, const char*, const char*, const char*, int, const char *, const char *, int);
extern void onValueGo(void*, const char*, valueHolder*);

void onDatumStart(void *rhPointer, const char *familyName, const char *groupName, const char *tagName, int typeId,
     const char *label, const char *interpretedValue, int numValues)
{
     onDatumStartGo(rhPointer, familyName, groupName, tagName, typeId, label, interpretedValue, numValues);
}

void onValue(void *rhPointer, const char *familyName, valueHolder *vh)
{
     onValueGo(rhPointer, familyName, vh);
}
*/
import "C"
