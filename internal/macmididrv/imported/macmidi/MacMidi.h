/**********************************************************************/
/*! \class RtMidi
    \brief An abstract base class for realtime MIDI input/output.

    This class implements some common functionality for the realtime
    MIDI input/output subclasses RtMidiIn and RtMidiOut.

    RtMidi GitHub site: https://github.com/thestk/rtmidi
    RtMidi WWW site: http://www.music.mcgill.ca/~gary/rtmidi/

    RtMidi: realtime MIDI i/o C++ classes
    Copyright (c) 2003-2023 Gary P. Scavone

    Permission is hereby granted, free of charge, to any person
    obtaining a copy of this software and associated documentation files
    (the "Software"), to deal in the Software without restriction,
    including without limitation the rights to use, copy, modify, merge,
    publish, distribute, sublicense, and/or sell copies of the Software,
    and to permit persons to whom the Software is furnished to do so,
    subject to the following conditions:

    The above copyright notice and this permission notice shall be
    included in all copies or substantial portions of the Software.

    Any person wishing to distribute modifications to the Software is
    asked to send the modifications to the original developer so that
    they can be incorporated into the canonical version.  This is,
    however, not a binding provision of this license.

    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
    EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
    MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
    IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR
    ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF
    CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
    WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
/**********************************************************************/

/*!
  \file RtMidi.h
 */

#ifndef MACMIDI_H
#define MACMIDI_H

#if defined _WIN32 || defined __CYGWIN__
#if defined(MACMIDI_EXPORT)
#define MACMIDI_DLL_PUBLIC __declspec(dllexport)
#else
#define MACMIDI_DLL_PUBLIC
#endif
#else
#if __GNUC__ >= 4
#define MACMIDI_DLL_PUBLIC __attribute__((visibility("default")))
#else
#define MACMIDI_DLL_PUBLIC
#endif
#endif

#define MACMIDI_VERSION_MAJOR 6
#define MACMIDI_VERSION_MINOR 0
#define MACMIDI_VERSION_PATCH 0
#define MACMIDI_VERSION_BETA 0

#define MACMIDI_TOSTRING2(n) #n
#define MACMIDI_TOSTRING(n) MACMIDI_TOSTRING2(n)

#if MACMIDI_VERSION_BETA > 0
#define MACMIDI_VERSION                                                        \
  MACMIDI_TOSTRING(MACMIDI_VERSION_MAJOR)                                      \
  "." MACMIDI_TOSTRING(MACMIDI_VERSION_MINOR) "." MACMIDI_TOSTRING(            \
      MACMIDI_VERSION_PATCH) "beta" MACMIDI_TOSTRING(MACMIDI_VERSION_BETA)
#else
#define MACMIDI_VERSION                                                        \
  MACMIDI_TOSTRING(MACMIDI_VERSION_MAJOR)                                      \
  "." MACMIDI_TOSTRING(MACMIDI_VERSION_MINOR) "." MACMIDI_TOSTRING(            \
      MACMIDI_VERSION_PATCH)
#endif

#include <CoreMIDI/CoreMIDI.h>
#include <exception>
#include <iostream>
#include <string>
#include <vector>

/************************************************************************/
/*! \class MacMidiError
    \brief Exception handling class for MacMidi.

    The MacMidiError class is quite simple but it does allow errors to be
    "caught" by MacMidiError::Type. See the MacMidi documentation to know
    which methods can throw an MacMidiError.
*/
/************************************************************************/

class MACMIDI_DLL_PUBLIC MacMidiError : public std::exception {
public:
  //! Defined MacMidiError types.
  enum Type {
    WARNING,           /*!< A non-critical error. */
    DEBUG_WARNING,     /*!< A non-critical error which might be useful for
                          debugging. */
    UNSPECIFIED,       /*!< The default, unspecified error type. */
    NO_DEVICES_FOUND,  /*!< No devices found on system. */
    INVALID_DEVICE,    /*!< An invalid device ID was specified. */
    MEMORY_ERROR,      /*!< An error occurred during memory allocation. */
    INVALID_PARAMETER, /*!< An invalid parameter was specified to a function. */
    INVALID_USE,       /*!< The function was called incorrectly. */
    DRIVER_ERROR,      /*!< A system driver error occurred. */
    SYSTEM_ERROR,      /*!< A system error occurred. */
    THREAD_ERROR       /*!< A thread error occurred. */
  };

  //! The constructor.
  MacMidiError(const std::string &message,
               Type type = MacMidiError::UNSPECIFIED) throw()
      : message_(message), type_(type) {}

  //! The destructor.
  virtual ~MacMidiError(void) throw() {}

  //! Prints thrown error message to stderr.
  virtual void printMessage(void) const throw() {
    std::cerr << '\n' << message_ << "\n\n";
  }

  //! Returns the thrown error message type.
  virtual const Type &getType(void) const throw() { return type_; }

  //! Returns the thrown error message string.
  virtual const std::string &getMessage(void) const throw() { return message_; }

  //! Returns the thrown error message as a c-style string.
  virtual const char *what(void) const throw() { return message_.c_str(); }

protected:
  std::string message_;
  Type type_;
};

//! MacMidi error callback function prototype.
/*!
    \param type Type of error.
    \param errorText Error description.

    Note that class behaviour is undefined after a critical error (not
    a warning) is reported.
 */
typedef void (*MacMidiErrorCallback)(MacMidiError::Type type,
                                     const std::string &errorText,
                                     void *userData);

void createCoreMidiClientSingleton(const std::string &clientName,
                                   MIDINotifyProc callback) throw();

class MidiApi;

class MACMIDI_DLL_PUBLIC MacMidi {
public:
  MacMidi(MacMidi &&other) noexcept;
  //! MIDI API specifier arguments.
  enum Api {
    UNSPECIFIED,   /*!< Search for a working compiled API. */
    MACOSX_CORE,   /*!< Macintosh OS-X CoreMIDI API. */
    LINUX_ALSA,    /*!< The Advanced Linux Sound Architecture API. */
    UNIX_JACK,     /*!< The JACK Low-Latency MIDI Server API. */
    WINDOWS_MM,    /*!< The Microsoft Multimedia MIDI API. */
    MACMIDI_DUMMY, /*!< A compilable but non-functional API. */
    WEB_MIDI_API,  /*!< W3C Web MIDI API. */
    WINDOWS_UWP,   /*!< The Microsoft Universal Windows Platform MIDI API. */
    ANDROID_AMIDI, /*!< Native Android MIDI API. */
    NUM_APIS       /*!< Number of values in this enum. */
  };

  //! A static function to determine the current MacMidi version.
  static std::string getVersion(void) throw();

  //! A static function to determine the available compiled MIDI APIs.
  /*!
    The values returned in the std::vector can be compared against
    the enumerated list values.  Note that there can be more than one
    API compiled for certain operating systems.
  */
  static void getCompiledApi(std::vector<MacMidi::Api> &apis) throw();

  //! Return the name of a specified compiled MIDI API.
  /*!
    This obtains a short lower-case name used for identification purposes.
    This value is guaranteed to remain identical across library versions.
    If the API is unknown, this function will return the empty string.
  */
  static std::string getApiName(MacMidi::Api api);

  //! Return the display name of a specified compiled MIDI API.
  /*!
    This obtains a long name used for display purposes.
    If the API is unknown, this function will return the empty string.
  */
  static std::string getApiDisplayName(MacMidi::Api api);

  //! Return the compiled MIDI API having the given name.
  /*!
    A case insensitive comparison will check the specified name
    against the list of compiled APIs, and return the one which
    matches. On failure, the function returns UNSPECIFIED.
  */
  static MacMidi::Api getCompiledApiByName(const std::string &name);

  //! Pure virtual openPort() function.
  virtual void
  openPort(unsigned int portNumber = 0,
           const std::string &portName = std::string("MacMidi")) = 0;

  //! Pure virtual openVirtualPort() function.
  virtual void
  openVirtualPort(const std::string &portName = std::string("MacMidi")) = 0;

  //! Pure virtual getPortCount() function.
  virtual unsigned int getPortCount() = 0;

  //! Pure virtual getPortName() function.
  virtual std::string getPortName(unsigned int portNumber = 0) = 0;

  //! Pure virtual closePort() function.
  virtual void closePort(void) = 0;

  void setClientName(const std::string &clientName);
  void setPortName(const std::string &portName);

  //! Returns true if a port is open and false if not.
  /*!
      Note that this only applies to connections made with the openPort()
      function, not to virtual ports.
  */
  virtual bool isPortOpen(void) const = 0;

  //! Set an error callback function to be invoked when an error has occurred.
  /*!
    The callback function will be called whenever an error has occurred. It is
    best to set the error callback function before opening a port.
  */
  virtual void setErrorCallback(MacMidiErrorCallback errorCallback = NULL,
                                void *userData = 0) = 0;

protected:
  MacMidi();
  virtual ~MacMidi();
  MidiApi *macapi_;

  /* Make the class non-copyable */
  MacMidi(MacMidi &other) = delete;
  MacMidi &operator=(MacMidi &other) = delete;
};

/**********************************************************************/
/*! \class MacMidiIn
    \brief A realtime MIDI input class.

    This class provides a common, platform-independent API for
    realtime MIDI input.  It allows access to a single MIDI input
    port.  Incoming MIDI messages are either saved to a queue for
    retrieval using the getMessage() function or immediately passed to
    a user-specified callback function.  Create multiple instances of
    this class to connect to more than one MIDI device at the same
    time.  With the OS-X, Linux ALSA, and JACK MIDI APIs, it is also
    possible to open a virtual input port to which other MIDI software
    clients can connect.
*/
/**********************************************************************/

// **************************************************************** //
//
// MacMidiIn and MacMidiOut class declarations.
//
// MacMidiIn / MacMidiOut are "controllers" used to select an available
// MIDI input or output interface.  They present common APIs for the
// user to call but all functionality is implemented by the classes
// MidiInApi, MidiOutApi and their subclasses.  MacMidiIn and MacMidiOut
// each create an instance of a MidiInApi or MidiOutApi subclass based
// on the user's API choice.  If no choice is made, they attempt to
// make a "logical" API selection.
//
// **************************************************************** //

class MACMIDI_DLL_PUBLIC MacMidiIn : public MacMidi {
public:
  //! User callback function type definition.
  typedef void (*MacMidiCallback)(double timeStamp,
                                  std::vector<unsigned char> *message,
                                  void *userData);

  //! Default constructor that allows an optional api, client name and queue
  //! size.
  /*!
    An exception will be thrown if a MIDI system initialization
    error occurs.  The queue size defines the maximum number of
    messages that can be held in the MIDI queue (when not using a
    callback function).  If the queue size limit is reached,
    incoming messages will be ignored.

    If no API argument is specified and multiple API support has been
    compiled, the default order of use is ALSA, JACK (Linux) and CORE,
    JACK (OS-X).

    \param api        An optional API id can be specified.
    \param clientName An optional client name can be specified. This
                      will be used to group the ports that are created
                      by the application.
    \param queueSizeLimit An optional size of the MIDI input queue can be
    specified.
  */
  MacMidiIn(MacMidi::Api api = UNSPECIFIED,
            const std::string &clientName = "MacMidi Input Client",
            unsigned int queueSizeLimit = 100);

  MacMidiIn(MacMidiIn &&other) noexcept : MacMidi(std::move(other)) {}

  //! If a MIDI connection is still open, it will be closed by the destructor.
  ~MacMidiIn(void) throw();

  //! Returns the MIDI API specifier for the current instance of MacMidiIn.
  MacMidi::Api getCurrentApi(void) throw();

  //! Open a MIDI input connection given by enumeration number.
  /*!
    \param portNumber An optional port number greater than 0 can be specified.
                      Otherwise, the default or first port found is opened.
    \param portName An optional name for the application port that is used to
    connect to portId can be specified.
  */
  void openPort(unsigned int portNumber = 0,
                const std::string &portName = std::string("MacMidi Input"));

  //! Create a virtual input port, with optional name, to allow software
  //! connections (OS X, JACK and ALSA only).
  /*!
    This function creates a virtual MIDI input port to which other
    software applications can connect.  This type of functionality
    is currently only supported by the Macintosh OS-X, any JACK,
    and Linux ALSA APIs (the function returns an error for the other APIs).

    \param portName An optional name for the application port that is
                    used to connect to portId can be specified.
  */
  void
  openVirtualPort(const std::string &portName = std::string("MacMidi Input"));

  //! Set a callback function to be invoked for incoming MIDI messages.
  /*!
    The callback function will be called whenever an incoming MIDI
    message is received.  While not absolutely necessary, it is best
    to set the callback function before opening a MIDI port to avoid
    leaving some messages in the queue.

    \param callback A callback function must be given.
    \param userData Optionally, a pointer to additional data can be
                    passed to the callback function whenever it is called.
  */
  void setCallback(MacMidiCallback callback, void *userData = 0);

  //! Cancel use of the current callback function (if one exists).
  /*!
    Subsequent incoming MIDI messages will be written to the queue
    and can be retrieved with the \e getMessage function.
  */
  void cancelCallback();

  //! Close an open MIDI connection (if one exists).
  void closePort(void);

  //! Returns true if a port is open and false if not.
  /*!
      Note that this only applies to connections made with the openPort()
      function, not to virtual ports.
  */
  virtual bool isPortOpen() const;

  //! Return the number of available MIDI input ports.
  /*!
    \return This function returns the number of MIDI ports of the selected API.
  */
  unsigned int getPortCount();

  //! Return a string identifier for the specified MIDI input port number.
  /*!
    \return The name of the port with the given Id is returned.
    \retval An empty string is returned if an invalid port specifier
            is provided. User code should assume a UTF-8 encoding.
  */
  std::string getPortName(unsigned int portNumber = 0);

  //! Specify whether certain MIDI message types should be queued or ignored
  //! during input.
  /*!
    By default, MIDI timing and active sensing messages are ignored
    during message input because of their relative high data rates.
    MIDI sysex messages are ignored by default as well.  Variable
    values of "true" imply that the respective message type will be
    ignored.
  */
  void ignoreTypes(bool midiSysex = true, bool midiTime = true,
                   bool midiSense = true);

  //! Fill the user-provided vector with the data bytes for the next available
  //! MIDI message in the input queue and return the event delta-time in
  //! seconds.
  /*!
    This function returns immediately whether a new message is
    available or not.  A valid message is indicated by a non-zero
    vector size.  An exception is thrown if an error occurs during
    message retrieval or an input connection was not previously
    established.
  */
  double getMessage(std::vector<unsigned char> *message);

  //! Set an error callback function to be invoked when an error has occurred.
  /*!
    The callback function will be called whenever an error has occurred. It is
    best to set the error callback function before opening a port.
  */
  virtual void setErrorCallback(MacMidiErrorCallback errorCallback = NULL,
                                void *userData = 0);

  //! Set maximum expected incoming message size.
  /*!
    For APIs that require manual buffer management, it can be useful to set the
    buffer size and buffer count when expecting to receive large SysEx messages.
    Note that currently this function has no effect when called after
    openPort().  The default buffer size is 1024 with a count of 4 buffers,
    which should be sufficient for most cases; as mentioned, this does not
    affect all API backends, since most either support dynamically scalable
    buffers or take care of buffer handling themselves.  It is principally
    intended for users of the Windows MM backend who must support receiving
    especially large messages.
  */
  virtual void setBufferSize(unsigned int size, unsigned int count);

protected:
  void openMidiApi(MacMidi::Api api, const std::string &clientName,
                   unsigned int queueSizeLimit);
};

/**********************************************************************/
/*! \class MacMidiOut
    \brief A realtime MIDI output class.

    This class provides a common, platform-independent API for MIDI
    output.  It allows one to probe available MIDI output ports, to
    connect to one such port, and to send MIDI bytes immediately over
    the connection.  Create multiple instances of this class to
    connect to more than one MIDI device at the same time.  With the
    OS-X, Linux ALSA and JACK MIDI APIs, it is also possible to open a
    virtual port to which other MIDI software clients can connect.
*/
/**********************************************************************/

class MACMIDI_DLL_PUBLIC MacMidiOut : public MacMidi {
public:
  //! Default constructor that allows an optional client name.
  /*!
    An exception will be thrown if a MIDI system initialization error occurs.

    If no API argument is specified and multiple API support has been
    compiled, the default order of use is ALSA, JACK (Linux) and CORE,
    JACK (OS-X).
  */
  MacMidiOut(MacMidi::Api api = UNSPECIFIED,
             const std::string &clientName = "MacMidi Output Client");

  MacMidiOut(MacMidiOut &&other) noexcept : MacMidi(std::move(other)) {}

  //! The destructor closes any open MIDI connections.
  ~MacMidiOut(void) throw();

  //! Returns the MIDI API specifier for the current instance of MacMidiOut.
  MacMidi::Api getCurrentApi(void) throw();

  //! Open a MIDI output connection.
  /*!
      An optional port number greater than 0 can be specified.
      Otherwise, the default or first port found is opened.  An
      exception is thrown if an error occurs while attempting to make
      the port connection.
  */
  void openPort(unsigned int portNumber = 0,
                const std::string &portName = std::string("MacMidi Output"));

  //! Close an open MIDI connection (if one exists).
  void closePort(void);

  //! Returns true if a port is open and false if not.
  /*!
      Note that this only applies to connections made with the openPort()
      function, not to virtual ports.
  */
  virtual bool isPortOpen() const;

  //! Create a virtual output port, with optional name, to allow software
  //! connections (OS X, JACK and ALSA only).
  /*!
      This function creates a virtual MIDI output port to which other
      software applications can connect.  This type of functionality
      is currently only supported by the Macintosh OS-X, Linux ALSA
      and JACK APIs (the function does nothing with the other APIs).
      An exception is thrown if an error occurs while attempting to
      create the virtual port.
  */
  void
  openVirtualPort(const std::string &portName = std::string("MacMidi Output"));

  //! Return the number of available MIDI output ports.
  unsigned int getPortCount(void);

  //! Return a string identifier for the specified MIDI port type and number.
  /*!
    \return The name of the port with the given Id is returned.
    \retval An empty string is returned if an invalid port specifier
            is provided. User code should assume a UTF-8 encoding.
  */
  std::string getPortName(unsigned int portNumber = 0);

  //! Immediately send a single message out an open MIDI output port.
  /*!
      An exception is thrown if an error occurs during output or an
      output connection was not previously established.
  */
  void sendMessage(const std::vector<unsigned char> *message);

  //! Immediately send a single message out an open MIDI output port.
  /*!
      An exception is thrown if an error occurs during output or an
      output connection was not previously established.

      \param message A pointer to the MIDI message as raw bytes
      \param size    Length of the MIDI message in bytes
  */
  void sendMessage(const unsigned char *message, size_t size);

  //! Set an error callback function to be invoked when an error has occurred.
  /*!
    The callback function will be called whenever an error has occurred. It is
    best to set the error callback function before opening a port.
  */
  virtual void setErrorCallback(MacMidiErrorCallback errorCallback = NULL,
                                void *userData = 0);

protected:
  void openMidiApi(MacMidi::Api api, const std::string &clientName);
};

// **************************************************************** //
//
// MidiInApi / MidiOutApi class declarations.
//
// Subclasses of MidiInApi and MidiOutApi contain all API- and
// OS-specific code necessary to fully implement the MacMidi API.
//
// Note that MidiInApi and MidiOutApi are abstract base classes and
// cannot be explicitly instantiated.  MacMidiIn and MacMidiOut will
// create instances of a MidiInApi or MidiOutApi subclass.
//
// **************************************************************** //

class MACMIDI_DLL_PUBLIC MidiApi {
public:
  MidiApi();
  virtual ~MidiApi();
  virtual MacMidi::Api getCurrentApi(void) = 0;
  virtual void openPort(unsigned int portNumber,
                        const std::string &portName) = 0;
  virtual void openVirtualPort(const std::string &portName) = 0;
  virtual void closePort(void) = 0;
  virtual void setClientName(const std::string &clientName) = 0;
  virtual void setPortName(const std::string &portName) = 0;

  virtual unsigned int getPortCount(void) = 0;
  virtual std::string getPortName(unsigned int portNumber) = 0;

  inline bool isPortOpen() const { return connected_; }
  void setErrorCallback(MacMidiErrorCallback errorCallback, void *userData);

  //! A basic error reporting function for MacMidi classes.
  void error(MacMidiError::Type type, std::string errorString);

protected:
  virtual void initialize(const std::string &clientName) = 0;

  void *apiData_;
  bool connected_;
  std::string errorString_;
  MacMidiErrorCallback errorCallback_;
  bool firstErrorOccurred_;
  void *errorCallbackUserData_;
};

class MACMIDI_DLL_PUBLIC MidiInApi : public MidiApi {
public:
  MidiInApi(unsigned int queueSizeLimit);
  virtual ~MidiInApi(void);
  void setCallback(MacMidiIn::MacMidiCallback callback, void *userData);
  void cancelCallback(void);
  virtual void ignoreTypes(bool midiSysex, bool midiTime, bool midiSense);
  virtual double getMessage(std::vector<unsigned char> *message);
  virtual void setBufferSize(unsigned int size, unsigned int count);

  // A MIDI structure used internally by the class to store incoming
  // messages.  Each message represents one and only one MIDI message.
  struct MidiMessage {
    std::vector<unsigned char> bytes;

    //! Time in seconds elapsed since the previous message
    double timeStamp;

    // Default constructor.
    MidiMessage() : bytes(0), timeStamp(0.0) {}
  };

  struct MidiQueue {
    unsigned int front;
    unsigned int back;
    unsigned int ringSize;
    MidiMessage *ring;

    // Default constructor.
    MidiQueue() : front(0), back(0), ringSize(0), ring(0) {}
    bool push(const MidiMessage &);
    bool pop(std::vector<unsigned char> *, double *);
    unsigned int size(unsigned int *back = 0, unsigned int *front = 0);
  };

  // The MacMidiInData structure is used to pass private class data to
  // the MIDI input handling function or thread.
  struct MacMidiInData {
    MidiQueue queue;
    MidiMessage message;
    unsigned char ignoreFlags;
    bool doInput;
    bool firstMessage;
    void *apiData;
    bool usingCallback;
    MacMidiIn::MacMidiCallback userCallback;
    void *userData;
    bool continueSysex;
    unsigned int bufferSize;
    unsigned int bufferCount;

    // Default constructor.
    MacMidiInData()
        : ignoreFlags(7), doInput(false), firstMessage(true), apiData(0),
          usingCallback(false), userCallback(0), userData(0),
          continueSysex(false), bufferSize(1024), bufferCount(4) {}
  };

protected:
  MacMidiInData inputData_;
};

class MACMIDI_DLL_PUBLIC MidiOutApi : public MidiApi {
public:
  MidiOutApi(void);
  virtual ~MidiOutApi(void);
  virtual void sendMessage(const unsigned char *message, size_t size) = 0;
};

// **************************************************************** //
//
// Inline MacMidiIn and MacMidiOut definitions.
//
// **************************************************************** //

inline MacMidi::Api MacMidiIn ::getCurrentApi(void) throw() {
  return macapi_->getCurrentApi();
}
inline void MacMidiIn ::openPort(unsigned int portNumber,
                                 const std::string &portName) {
  macapi_->openPort(portNumber, portName);
}
inline void MacMidiIn ::openVirtualPort(const std::string &portName) {
  macapi_->openVirtualPort(portName);
}
inline void MacMidiIn ::closePort(void) { macapi_->closePort(); }
inline bool MacMidiIn ::isPortOpen() const { return macapi_->isPortOpen(); }
inline void MacMidiIn ::setCallback(MacMidiCallback callback, void *userData) {
  static_cast<MidiInApi *>(macapi_)->setCallback(callback, userData);
}
inline void MacMidiIn ::cancelCallback(void) {
  static_cast<MidiInApi *>(macapi_)->cancelCallback();
}
inline unsigned int MacMidiIn ::getPortCount(void) {
  return macapi_->getPortCount();
}
inline std::string MacMidiIn ::getPortName(unsigned int portNumber) {
  return macapi_->getPortName(portNumber);
}
inline void MacMidiIn ::ignoreTypes(bool midiSysex, bool midiTime,
                                    bool midiSense) {
  static_cast<MidiInApi *>(macapi_)->ignoreTypes(midiSysex, midiTime,
                                                 midiSense);
}
inline double MacMidiIn ::getMessage(std::vector<unsigned char> *message) {
  return static_cast<MidiInApi *>(macapi_)->getMessage(message);
}
inline void MacMidiIn ::setErrorCallback(MacMidiErrorCallback errorCallback,
                                         void *userData) {
  macapi_->setErrorCallback(errorCallback, userData);
}
inline void MacMidiIn ::setBufferSize(unsigned int size, unsigned int count) {
  static_cast<MidiInApi *>(macapi_)->setBufferSize(size, count);
}

inline MacMidi::Api MacMidiOut ::getCurrentApi(void) throw() {
  return macapi_->getCurrentApi();
}
inline void MacMidiOut ::openPort(unsigned int portNumber,
                                  const std::string &portName) {
  macapi_->openPort(portNumber, portName);
}
inline void MacMidiOut ::openVirtualPort(const std::string &portName) {
  macapi_->openVirtualPort(portName);
}
inline void MacMidiOut ::closePort(void) { macapi_->closePort(); }
inline bool MacMidiOut ::isPortOpen() const { return macapi_->isPortOpen(); }
inline unsigned int MacMidiOut ::getPortCount(void) {
  return macapi_->getPortCount();
}
inline std::string MacMidiOut ::getPortName(unsigned int portNumber) {
  return macapi_->getPortName(portNumber);
}
inline void
MacMidiOut ::sendMessage(const std::vector<unsigned char> *message) {
  static_cast<MidiOutApi *>(macapi_)->sendMessage(&message->at(0),
                                                  message->size());
}
inline void MacMidiOut ::sendMessage(const unsigned char *message,
                                     size_t size) {
  static_cast<MidiOutApi *>(macapi_)->sendMessage(message, size);
}
inline void MacMidiOut ::setErrorCallback(MacMidiErrorCallback errorCallback,
                                          void *userData) {
  macapi_->setErrorCallback(errorCallback, userData);
}

#endif
