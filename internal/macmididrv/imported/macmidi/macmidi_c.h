/************************************************************************/
/*! \defgroup C-interface
    @{

    \brief C interface to realtime MIDI input/output C++ classes.

    RtMidi offers a C-style interface, principally for use in binding
    RtMidi to other programming languages.  All structs, enums, and
    functions listed here have direct analogs (and simply call to)
    items in the C++ RtMidi class and its supporting classes and
    types
*/
/************************************************************************/

/*!
  \file rtmidi_c.h
 */

#include <stdbool.h>
#include <stddef.h>
#ifndef MACMIDI_C_H
#define MACMIDI_C_H

#if defined(MACMIDI_EXPORT)
#if defined _WIN32 || defined __CYGWIN__
#define MACMIDIAPI __declspec(dllexport)
#else
#define MACMIDIAPI __attribute__((visibility("default")))
#endif
#else
#define MACMIDIAPI //__declspec(dllimport)
#endif

#ifdef __cplusplus
extern "C" {
#endif

//! \brief Wraps an MacMidi object for C function return statuses.
struct MacMidiWrapper {
  //! The wrapped MacMidi object.
  void *ptr;
  void *data;

  //! True when the last function call was OK.
  bool ok;

  //! If an error occurred (ok != true), set to an error message.
  const char *msg;
};

//! \brief Typedef for a generic MacMidi pointer.
typedef struct MacMidiWrapper *MacMidiPtr;

//! \brief Typedef for a generic MacMidiIn pointer.
typedef struct MacMidiWrapper *MacMidiInPtr;

//! \brief Typedef for a generic MacMidiOut pointer.
typedef struct MacMidiWrapper *MacMidiOutPtr;

//! \brief MIDI API specifier arguments.  See \ref MacMidi::Api.
enum MacMidiApi {
  MACMIDI_API_UNSPECIFIED,   /*!< Search for a working compiled API. */
  MACMIDI_API_MACOSX_CORE,   /*!< Macintosh OS-X CoreMIDI API. */
  MACMIDI_API_LINUX_ALSA,    /*!< The Advanced Linux Sound Architecture API. */
  MACMIDI_API_UNIX_JACK,     /*!< The Jack Low-Latency MIDI Server API. */
  MACMIDI_API_WINDOWS_MM,    /*!< The Microsoft Multimedia MIDI API. */
  MACMIDI_API_MACMIDI_DUMMY, /*!< A compilable but non-functional API. */
  MACMIDI_API_WEB_MIDI_API,  /*!< W3C Web MIDI API. */
  MACMIDI_API_WINDOWS_UWP,   /*!< The Microsoft Universal Windows Platform MIDI
                                API. */
  MACMIDI_API_ANDROID,       /*!< The Android MIDI API. */
  MACMIDI_API_NUM            /*!< Number of values in this enum. */
};

//! \brief Defined MacMidiError types. See \ref MacMidiError::Type.
enum MacMidiErrorType {
  MACMIDI_ERROR_WARNING,       /*!< A non-critical error. */
  MACMIDI_ERROR_DEBUG_WARNING, /*!< A non-critical error which might be useful
                                  for debugging. */
  MACMIDI_ERROR_UNSPECIFIED,   /*!< The default, unspecified error type. */
  MACMIDI_ERROR_NO_DEVICES_FOUND, /*!< No devices found on system. */
  MACMIDI_ERROR_INVALID_DEVICE,   /*!< An invalid device ID was specified. */
  MACMIDI_ERROR_MEMORY_ERROR, /*!< An error occurred during memory allocation.
                               */
  MACMIDI_ERROR_INVALID_PARAMETER, /*!< An invalid parameter was specified to a
                                      function. */
  MACMIDI_ERROR_INVALID_USE,       /*!< The function was called incorrectly. */
  MACMIDI_ERROR_DRIVER_ERROR,      /*!< A system driver error occurred. */
  MACMIDI_ERROR_SYSTEM_ERROR,      /*!< A system error occurred. */
  MACMIDI_ERROR_THREAD_ERROR       /*!< A thread error occurred. */
};

/*! \brief The type of a MacMidi callback function.
 *
 * \param timeStamp   The time at which the message has been received.
 * \param message     The midi message.
 * \param userData    Additional user data for the callback.
 *
 * See \ref MacMidiIn::MacMidiCallback.
 */
typedef void (*MacMidiCCallback)(double timeStamp, const unsigned char *message,
                                 size_t messageSize, void *userData);

/* MacMidi API */

/*! \brief Return the current MacMidi version.
 *! See \ref MacMidi::getVersion().
 */
MACMIDIAPI const char *macmidi_get_version();

/*! \brief Determine the available compiled MIDI APIs.
 *
 * If the given `apis` parameter is null, returns the number of available APIs.
 * Otherwise, fill the given apis array with the MacMidi::Api values.
 *
 * \param apis  An array or a null value.
 * \param apis_size  Number of elements pointed to by apis
 * \return number of items needed for apis array if apis==NULL, or
 *         number of items written to apis array otherwise.  A negative
 *         return value indicates an error.
 *
 * See \ref MacMidi::getCompiledApi().
 */
MACMIDIAPI int macmidi_get_compiled_api(enum MacMidiApi *apis,
                                        unsigned int apis_size);

//! \brief Return the name of a specified compiled MIDI API.
//! See \ref MacMidi::getApiName().
MACMIDIAPI const char *macmidi_api_name(enum MacMidiApi api);

//! \brief Return the display name of a specified compiled MIDI API.
//! See \ref MacMidi::getApiDisplayName().
MACMIDIAPI const char *macmidi_api_display_name(enum MacMidiApi api);

//! \brief Return the compiled MIDI API having the given name.
//! See \ref MacMidi::getCompiledApiByName().
MACMIDIAPI enum MacMidiApi macmidi_compiled_api_by_name(const char *name);

//! \internal Report an error.
MACMIDIAPI void macmidi_error(enum MacMidiErrorType type,
                              const char *errorString);

/*! \brief Open a MIDI port.
 *
 * \param port      Must be greater than 0
 * \param portName  Name for the application port.
 *
 * See MacMidi::openPort().
 */
MACMIDIAPI void macmidi_open_port(MacMidiPtr device, unsigned int portNumber,
                                  const char *portName);

/*! \brief Creates a virtual MIDI port to which other software applications can
 * connect.
 *
 * \param portName  Name for the application port.
 *
 * See MacMidi::openVirtualPort().
 */
MACMIDIAPI void macmidi_open_virtual_port(MacMidiPtr device,
                                          const char *portName);

/*! \brief Close a MIDI connection.
 * See MacMidi::closePort().
 */
MACMIDIAPI void macmidi_close_port(MacMidiPtr device);

/*! \brief Return the number of available MIDI ports.
 * See MacMidi::getPortCount().
 */
MACMIDIAPI unsigned int macmidi_get_port_count(MacMidiPtr device);

/*! \brief Access a string identifier for the specified MIDI input port number.
 *
 * To prevent memory leaks a char buffer must be passed to this function.
 * NULL can be passed as bufOut parameter, and that will write the required
 * buffer length in the bufLen.
 *
 * See MacMidi::getPortName().
 */
MACMIDIAPI int macmidi_get_port_name(MacMidiPtr device, unsigned int portNumber,
                                     char *bufOut, int *bufLen);

/* MacMidiIn API */

//! \brief Create a default MacMidiInPtr value, with no initialization.
MACMIDIAPI MacMidiInPtr macmidi_in_create_default(void);

/*! \brief Create a  MacMidiInPtr value, with given api, clientName and
 * queueSizeLimit.
 *
 *  \param api            An optional API id can be specified.
 *  \param clientName     An optional client name can be specified. This
 *                        will be used to group the ports that are created
 *                        by the application.
 *  \param queueSizeLimit An optional size of the MIDI input queue can be
 *                        specified.
 *
 * See MacMidiIn::MacMidiIn().
 */
MACMIDIAPI MacMidiInPtr macmidi_in_create(enum MacMidiApi api,
                                          const char *clientName,
                                          unsigned int queueSizeLimit);

//! \brief Free the given MacMidiInPtr.
MACMIDIAPI void macmidi_in_free(MacMidiInPtr device);

//! \brief Returns the MIDI API specifier for the given instance of MacMidiIn.
//! See \ref MacMidiIn::getCurrentApi().
MACMIDIAPI enum MacMidiApi macmidi_in_get_current_api(MacMidiPtr device);

//! \brief Set a callback function to be invoked for incoming MIDI messages.
//! See \ref MacMidiIn::setCallback().
MACMIDIAPI void macmidi_in_set_callback(MacMidiInPtr device,
                                        MacMidiCCallback callback,
                                        void *userData);

//! \brief Cancel use of the current callback function (if one exists).
//! See \ref MacMidiIn::cancelCallback().
MACMIDIAPI void macmidi_in_cancel_callback(MacMidiInPtr device);

//! \brief Specify whether certain MIDI message types should be queued or
//! ignored during input. See \ref MacMidiIn::ignoreTypes().
MACMIDIAPI void macmidi_in_ignore_types(MacMidiInPtr device, bool midiSysex,
                                        bool midiTime, bool midiSense);

/*! Fill the user-provided array with the data bytes for the next available
 * MIDI message in the input queue and return the event delta-time in seconds.
 *
 * \param message   Must point to a char* that is already allocated.
 *                  SYSEX messages maximum size being 1024, a statically
 *                  allocated array could
 *                  be sufficient.
 * \param size      Is used to return the size of the message obtained.
 *                  Must be set to the size of \ref message when calling.
 *
 * See MacMidiIn::getMessage().
 */
MACMIDIAPI double macmidi_in_get_message(MacMidiInPtr device,
                                         unsigned char *message, size_t *size);

/* MacMidiOut API */

//! \brief Create a default MacMidiInPtr value, with no initialization.
MACMIDIAPI MacMidiOutPtr macmidi_out_create_default(void);

/*! \brief Create a MacMidiOutPtr value, with given and clientName.
 *
 *  \param api            An optional API id can be specified.
 *  \param clientName     An optional client name can be specified. This
 *                        will be used to group the ports that are created
 *                        by the application.
 *
 * See MacMidiOut::MacMidiOut().
 */
MACMIDIAPI MacMidiOutPtr macmidi_out_create(enum MacMidiApi api,
                                            const char *clientName);

//! \brief Free the given MacMidiOutPtr.
MACMIDIAPI void macmidi_out_free(MacMidiOutPtr device);

//! \brief Returns the MIDI API specifier for the given instance of MacMidiOut.
//! See \ref MacMidiOut::getCurrentApi().
MACMIDIAPI enum MacMidiApi macmidi_out_get_current_api(MacMidiPtr device);

//! \brief Immediately send a single message out an open MIDI output port.
//! See \ref MacMidiOut::sendMessage().
MACMIDIAPI int macmidi_out_send_message(MacMidiOutPtr device,
                                        const unsigned char *message,
                                        int length);

#ifdef __cplusplus
}
#endif
#endif

/*! }@ */
