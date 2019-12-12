package ezif

/*
#include "exiv2_bridge.h"

extern void onDatumEndGo(void*, const char*);
extern void onDatumStartGo(void*, const char*, const char*, const char*, int, const char *, const char *, int, int);
extern void onValueGo(void*, valueHolder*);

void onDatumEnd(void *rhPointer, const char *familyName)
{
     onDatumEndGo(rhPointer, familyName);
}

void onDatumStart(void *rhPointer, const char *familyName, const char *groupName, const char *tagName, int typeId,
     const char *label, const char *interpretedValue, int numValues, int repeatable)
{
     onDatumStartGo(rhPointer, familyName, groupName, tagName, typeId, label, interpretedValue, numValues, repeatable);
}

void onValue(void *rhPointer, valueHolder *vh)
{
     onValueGo(rhPointer, vh);
}
*/
import "C"
