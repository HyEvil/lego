#include "hydra.h"
void CallHydraResume(HydraResume f, void* goEx, unsigned int coroutineId, char suc, void* data, unsigned int size)
{
    f(goEx, coroutineId, suc, data, size);
}

void CallHydraRet(HydraRet f, void* goEx, void* data, unsigned int size) { f(goEx, data, size); }