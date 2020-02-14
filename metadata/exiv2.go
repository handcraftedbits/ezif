package metadata // import "golang.handcraftedbits.com/ezif/metadata"

/*
#cgo LDFLAGS: -lexiv2 -lexiv2-xmp -lexpat -lz

#include "exiv2.h"

extern void onPropertyEndGo(void*, const char*);
extern void onPropertyStartGo(void*, const char*, const char*, const char*, int, const char *, const char *, int, int);
extern void onValueGo(void*, valueHolder*);

void onPropertyEnd(void *rhPointer, const char *familyName)
{
     onPropertyEndGo(rhPointer, familyName);
}

void onPropertyStart(void *rhPointer, const char *familyName, const char *groupName, const char *tagName, int typeId,
     const char *label, const char *interpretedValue, int numValues, int repeatable)
{
     onPropertyStartGo(rhPointer, familyName, groupName, tagName, typeId, label, interpretedValue, numValues,
          repeatable);
}

void onValue(void *rhPointer, valueHolder *vh)
{
     onValueGo(rhPointer, vh);
}
*/
import "C"
