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
     long longValue;
     int minuteValue;
     int monthValue;
     const char *strValue;
     uint32_t rationalValueD;
     uint32_t rationalValueN;
     int secondValue;
     int timezoneHourOffset;
     int timezoneMinuteOffset;
     int yearValue;
} valueHolder;

// Function pointer definitions

typedef void (*datumOnEndCallback)(void*, const char *);
typedef void (*datumOnStartCallback)(void*, const char *, const char *, const char *, int, const char *, const char *,
     int);
typedef void (*valueCallback)(void*, valueHolder*);

// Struct definitions

typedef struct readHandlers
{
     datumOnEndCallback doec;
     datumOnStartCallback dosc;
     valueCallback vc;
} readHandlers;

// Function definitions

void onDatumEnd(void*, const char*);
void onDatumStart(void*, const char*, const char*, const char*, int, const char *, const char *, int);
void onValue(void*, valueHolder*);
void readImageMetadata (const char*, exiv2Error*, valueHolder*, readHandlers*, void*);

#ifdef __cplusplus
}
#endif

#endif
