
#ifndef EFML_EXPORT_H
#define EFML_EXPORT_H

#ifdef EFML_STATIC_DEFINE
#  define EFML_EXPORT
#  define EFML_NO_EXPORT
#else
#  ifndef EFML_EXPORT
#    ifdef libefml_EXPORTS
        /* We are building this library */
#      define EFML_EXPORT __attribute__((visibility("default")))
#    else
        /* We are using this library */
#      define EFML_EXPORT __attribute__((visibility("default")))
#    endif
#  endif

#  ifndef EFML_NO_EXPORT
#    define EFML_NO_EXPORT __attribute__((visibility("hidden")))
#  endif
#endif

#ifndef EFML_DEPRECATED
#  define EFML_DEPRECATED __attribute__ ((__deprecated__))
#endif

#ifndef EFML_DEPRECATED_EXPORT
#  define EFML_DEPRECATED_EXPORT EFML_EXPORT EFML_DEPRECATED
#endif

#ifndef EFML_DEPRECATED_NO_EXPORT
#  define EFML_DEPRECATED_NO_EXPORT EFML_NO_EXPORT EFML_DEPRECATED
#endif

#if 0 /* DEFINE_NO_DEPRECATED */
#  ifndef EFML_NO_DEPRECATED
#    define EFML_NO_DEPRECATED
#  endif
#endif

#endif /* EFML_EXPORT_H */
