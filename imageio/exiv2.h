#ifndef _EXIV2_BRIDGE_H_
#define _EXIV2_BRIDGE_H_

#include <stdint.h>

#ifdef __cplusplus
extern "C"
{
#endif

// Struct definitions

typedef struct exiv2Error
{
     int code;
     const char *message;
} exiv2Error;

typedef struct valueHolder
{
     int dayValue;
     double doubleValue;
     int hourValue;
     const char *langValue;
     long longValue;
     int minuteValue;
     int monthValue;
     char *strValue;
     uint32_t rationalValueD;
     uint32_t rationalValueN;
     int secondValue;
     int timezoneHourOffset;
     int timezoneMinuteOffset;
     int yearValue;
} valueHolder;

// Function pointer definitions

typedef void (*propertyOnEndCallback)(void*, const char *);
typedef void (*propertyOnStartCallback)(void*, const char *, const char *, const char *, int, const char *, const char *,
     int, int);
typedef void (*valueCallback)(void*, valueHolder*);

// Struct definitions

typedef struct readHandler
{
     propertyOnEndCallback poec;
     propertyOnStartCallback posc;
     valueCallback vc;
} readHandler;

// Function definitions

void onPropertyEnd(void*, const char*);
void onPropertyStart(void*, const char*, const char*, const char*, int, const char *, const char *, int, int);
void onValue(void*, valueHolder*);
void readImageMetadata (const char*, exiv2Error*, valueHolder*, readHandler*, void*);

#ifdef __cplusplus
}
#endif

#endif
