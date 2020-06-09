#ifndef HYDRA_HEADER
#define HYDRA_HEADER

#include <stdbool.h>
typedef void (*HydraResume)(void* goEx, unsigned int coroutineId, char suc, const void* data, unsigned int size);
typedef void (*HydraRet)(void* goEx, void* data, unsigned int size);
void CallHydraResume(HydraResume f, void* goEx, unsigned int coroutineId, char suc, void* data, unsigned int size);
void CallHydraRet(HydraRet f, void* goEx, void* data, unsigned int size);
#endif