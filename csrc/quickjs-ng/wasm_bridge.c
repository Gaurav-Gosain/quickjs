#include "quickjs.h"
#include <string.h>
#include <stdint.h>

#define HEAP_SIZE (16 * 1024 * 1024)
static char heap[HEAP_SIZE];
static size_t heap_ptr = 0;

static void* my_calloc(void *opaque, size_t count, size_t size) {
    size_t total = (count * size + 7) & ~7;
    if (heap_ptr + total > HEAP_SIZE) return NULL;
    void* ptr = &heap[heap_ptr];
    memset(ptr, 0, total);
    heap_ptr += total;
    return ptr;
}

static void* my_malloc(void *opaque, size_t size) {
    size = (size + 7) & ~7;
    if (heap_ptr + size > HEAP_SIZE) return NULL;
    void* ptr = &heap[heap_ptr];
    heap_ptr += size;
    return ptr;
}

static void my_free(void *opaque, void *ptr) { (void)opaque; (void)ptr; }

static void* my_realloc(void *opaque, void *ptr, size_t size) {
    if (!ptr) return my_malloc(opaque, size);
    void* new_ptr = my_malloc(opaque, size);
    if (new_ptr && ptr) memcpy(new_ptr, ptr, size);
    return new_ptr;
}

static size_t my_malloc_usable_size(const void *ptr) { return 0; }

static const JSMallocFunctions my_mf = {
    my_calloc, my_malloc, my_free, my_realloc, my_malloc_usable_size,
};

__attribute__((export_name("qjs_new_runtime")))
uint32_t qjs_new_runtime(void) {
    return (uint32_t)(uintptr_t)JS_NewRuntime2(&my_mf, NULL);
}

__attribute__((export_name("qjs_free_runtime")))
void qjs_free_runtime(uint32_t rt) {
    JS_FreeRuntime((JSRuntime*)(uintptr_t)rt);
}

__attribute__((export_name("qjs_new_context")))
uint32_t qjs_new_context(uint32_t rt) {
    return (uint32_t)(uintptr_t)JS_NewContext((JSRuntime*)(uintptr_t)rt);
}

__attribute__((export_name("qjs_free_context")))
void qjs_free_context(uint32_t ctx) {
    JS_FreeContext((JSContext*)(uintptr_t)ctx);
}

// For JSValue, we store it in memory and return pointer
// This avoids 64-bit value issues in 32-bit WASM
__attribute__((export_name("qjs_eval")))
uint32_t qjs_eval(uint32_t ctx, uint32_t code_ptr, uint32_t len, uint32_t filename_ptr, int flags) {
    const char* code = (const char*)(uintptr_t)code_ptr;
    const char* filename = (const char*)(uintptr_t)filename_ptr;
    
    // Allocate space for JSValue result
    JSValue* result_ptr = (JSValue*)my_malloc(NULL, sizeof(JSValue));
    if (!result_ptr) return 0;
    
    *result_ptr = JS_Eval((JSContext*)(uintptr_t)ctx, code, len, filename, flags);
    return (uint32_t)(uintptr_t)result_ptr;
}

__attribute__((export_name("qjs_is_exception")))
int qjs_is_exception(uint32_t val_ptr) {
    if (!val_ptr) return 1;
    JSValue* vp = (JSValue*)(uintptr_t)val_ptr;
    return JS_IsException(*vp);
}

__attribute__((export_name("qjs_is_undefined")))
int qjs_is_undefined(uint32_t val_ptr) {
    if (!val_ptr) return 1;
    JSValue* vp = (JSValue*)(uintptr_t)val_ptr;
    return JS_IsUndefined(*vp);
}

__attribute__((export_name("qjs_to_cstring")))
uint32_t qjs_to_cstring(uint32_t ctx, uint32_t val_ptr) {
    if (!val_ptr) return 0;
    JSValue* vp = (JSValue*)(uintptr_t)val_ptr;
    return (uint32_t)(uintptr_t)JS_ToCString((JSContext*)(uintptr_t)ctx, *vp);
}

__attribute__((export_name("qjs_free_cstring")))
void qjs_free_cstring(uint32_t ctx, uint32_t str_ptr) {
    JS_FreeCString((JSContext*)(uintptr_t)ctx, (const char*)(uintptr_t)str_ptr);
}

__attribute__((export_name("qjs_free_value")))
void qjs_free_value(uint32_t ctx, uint32_t val_ptr) {
    if (!val_ptr) return;
    JSValue* vp = (JSValue*)(uintptr_t)val_ptr;
    JS_FreeValue((JSContext*)(uintptr_t)ctx, *vp);
}

__attribute__((export_name("qjs_get_exception")))
uint32_t qjs_get_exception(uint32_t ctx) {
    JSValue* result_ptr = (JSValue*)my_malloc(NULL, sizeof(JSValue));
    if (!result_ptr) return 0;
    *result_ptr = JS_GetException((JSContext*)(uintptr_t)ctx);
    return (uint32_t)(uintptr_t)result_ptr;
}

__attribute__((export_name("qjs_alloc")))
uint32_t qjs_alloc(uint32_t size) {
    return (uint32_t)(uintptr_t)my_malloc(NULL, size);
}

__attribute__((export_name("qjs_get_heap_ptr")))
uint32_t qjs_get_heap_ptr(void) {
    return (uint32_t)heap_ptr;
}

__attribute__((export_name("qjs_reset_heap")))
void qjs_reset_heap(void) {
    heap_ptr = 0;
}
