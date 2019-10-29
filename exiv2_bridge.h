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
     double doubleValue;
     long longValue;
     const char *strValue;
     uint32_t rationalValueD;
     uint32_t rationalValueN;
} valueHolder;

// Function pointer definitions

typedef void (*datumOnStartCallback)(void *, const char *, const char *, const char *, int, const char *, const char *,
     int);
typedef void (*valueCallback)(void*, const char*, valueHolder*);

// Struct definitions

typedef struct readHandlers
{
     datumOnStartCallback dosc;
     valueCallback vc;
} readHandlers;

// Function definitions

void onDatumStart(void*, const char*, const char*, const char*, int, const char *, const char *, int);
void onValue(void*, const char*, valueHolder*);
void readImageMetadata (const char*, exiv2Error*, valueHolder*, readHandlers*, void*);

#ifdef __cplusplus
}
#endif

#endif
