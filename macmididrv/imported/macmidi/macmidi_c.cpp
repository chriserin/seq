#include "macmidi_c.h"
#include "MacMidi.h"
#include <stdlib.h>
#include <string.h>

/* Compile-time assertions that will break if the enums are changed in
 * the future without synchronizing them properly.  If you get (g++)
 * "error: ‘StaticEnumAssert<b>::StaticEnumAssert() [with bool b = false]’
 * is private within this context", it means enums are not aligned. */
template <bool b> class StaticEnumAssert {
private:
  StaticEnumAssert() {}
};
template <> class StaticEnumAssert<true> {
public:
  StaticEnumAssert() {}
};
#define ENUM_EQUAL(x, y) StaticEnumAssert<(int)x == (int)y>()
class StaticEnumAssertions {
  StaticEnumAssertions() {
    ENUM_EQUAL(MACMIDI_API_UNSPECIFIED, MacMidi::UNSPECIFIED);
    ENUM_EQUAL(MACMIDI_API_MACOSX_CORE, MacMidi::MACOSX_CORE);
    ENUM_EQUAL(MACMIDI_API_LINUX_ALSA, MacMidi::LINUX_ALSA);
    ENUM_EQUAL(MACMIDI_API_UNIX_JACK, MacMidi::UNIX_JACK);
    ENUM_EQUAL(MACMIDI_API_WINDOWS_MM, MacMidi::WINDOWS_MM);
    ENUM_EQUAL(MACMIDI_API_ANDROID, MacMidi::ANDROID_AMIDI);
    ENUM_EQUAL(MACMIDI_API_MACMIDI_DUMMY, MacMidi::MACMIDI_DUMMY);
    ENUM_EQUAL(MACMIDI_API_WEB_MIDI_API, MacMidi::WEB_MIDI_API);
    ENUM_EQUAL(MACMIDI_API_WINDOWS_UWP, MacMidi::WINDOWS_UWP);

    ENUM_EQUAL(MACMIDI_ERROR_WARNING, MacMidiError::WARNING);
    ENUM_EQUAL(MACMIDI_ERROR_DEBUG_WARNING, MacMidiError::DEBUG_WARNING);
    ENUM_EQUAL(MACMIDI_ERROR_UNSPECIFIED, MacMidiError::UNSPECIFIED);
    ENUM_EQUAL(MACMIDI_ERROR_NO_DEVICES_FOUND, MacMidiError::NO_DEVICES_FOUND);
    ENUM_EQUAL(MACMIDI_ERROR_INVALID_DEVICE, MacMidiError::INVALID_DEVICE);
    ENUM_EQUAL(MACMIDI_ERROR_MEMORY_ERROR, MacMidiError::MEMORY_ERROR);
    ENUM_EQUAL(MACMIDI_ERROR_INVALID_PARAMETER,
               MacMidiError::INVALID_PARAMETER);
    ENUM_EQUAL(MACMIDI_ERROR_INVALID_USE, MacMidiError::INVALID_USE);
    ENUM_EQUAL(MACMIDI_ERROR_DRIVER_ERROR, MacMidiError::DRIVER_ERROR);
    ENUM_EQUAL(MACMIDI_ERROR_SYSTEM_ERROR, MacMidiError::SYSTEM_ERROR);
    ENUM_EQUAL(MACMIDI_ERROR_THREAD_ERROR, MacMidiError::THREAD_ERROR);
  }
};

class CallbackProxyUserData {
public:
  CallbackProxyUserData(MacMidiCCallback cCallback, void *userData)
      : c_callback(cCallback), user_data(userData) {}
  MacMidiCCallback c_callback;
  void *user_data;
};

#ifndef MACMIDI_SOURCE_INCLUDED
extern "C" const enum MacMidiApi
    macmidi_compiled_apis[]; // casting from MacMidi::Api[]
#endif
extern "C" const unsigned int macmidi_num_compiled_apis;

/* MacMidi API */
const char *macmidi_get_version() { return MACMIDI_VERSION; }

int macmidi_get_compiled_api(enum MacMidiApi *apis, unsigned int apis_size) {
  unsigned num = macmidi_num_compiled_apis;
  if (apis) {
    num = (num < apis_size) ? num : apis_size;
    memcpy(apis, macmidi_compiled_apis, num * sizeof(enum MacMidiApi));
  }
  return (int)num;
}

extern "C" const char *macmidi_api_names[][2];
const char *macmidi_api_name(enum MacMidiApi api) {
  if (api < 0 || api >= MACMIDI_API_NUM)
    return NULL;
  return macmidi_api_names[api][0];
}

const char *macmidi_api_display_name(enum MacMidiApi api) {
  if (api < 0 || api >= MACMIDI_API_NUM)
    return "Unknown";
  return macmidi_api_names[api][1];
}

enum MacMidiApi macmidi_compiled_api_by_name(const char *name) {
  MacMidi::Api api = MacMidi::UNSPECIFIED;
  if (name) {
    api = MacMidi::getCompiledApiByName(name);
  }
  return (enum MacMidiApi)api;
}

void macmidi_error(MidiApi *api, enum MacMidiErrorType type,
                   const char *errorString) {
  std::string msg = errorString;
  api->error((MacMidiError::Type)type, msg);
}

void macmidi_open_port(MacMidiPtr device, unsigned int portNumber,
                       const char *portName) {
  std::string name = portName;
  try {
    ((MacMidi *)device->ptr)->openPort(portNumber, name);

  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
  }
}

void macmidi_open_virtual_port(MacMidiPtr device, const char *portName) {
  std::string name = portName;
  try {
    ((MacMidi *)device->ptr)->openVirtualPort(name);

  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
  }
}

void macmidi_close_port(MacMidiPtr device) {
  try {
    ((MacMidi *)device->ptr)->closePort();

  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
  }
}

unsigned int macmidi_get_port_count(MacMidiPtr device) {
  try {
    return ((MacMidi *)device->ptr)->getPortCount();

  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
    return -1;
  }
}

int macmidi_get_port_name(MacMidiPtr device, unsigned int portNumber,
                          char *bufOut, int *bufLen) {
  if (bufOut == nullptr && bufLen == nullptr) {
    return -1;
  }

  std::string name;
  try {
    name = ((MacMidi *)device->ptr)->getPortName(portNumber);
  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
    return -1;
  }

  if (bufOut == nullptr) {
    *bufLen = static_cast<int>(name.size()) + 1;
    return 0;
  }

  return snprintf(bufOut, static_cast<size_t>(*bufLen), "%s", name.c_str());
}

/* MacMidiIn API */
MacMidiInPtr macmidi_in_create_default() {
  MacMidiWrapper *wrp = new MacMidiWrapper;

  try {
    MacMidiIn *rIn = new MacMidiIn();

    wrp->ptr = (void *)rIn;
    wrp->data = 0;
    wrp->ok = true;
    wrp->msg = "";

  } catch (const MacMidiError &err) {
    wrp->ptr = 0;
    wrp->data = 0;
    wrp->ok = false;
    wrp->msg = err.what();
  }

  return wrp;
}

MacMidiInPtr macmidi_in_create(enum MacMidiApi api, const char *clientName,
                               unsigned int queueSizeLimit) {
  std::string name = clientName;
  MacMidiWrapper *wrp = new MacMidiWrapper;

  try {
    MacMidiIn *rIn = new MacMidiIn((MacMidi::Api)api, name, queueSizeLimit);

    wrp->ptr = (void *)rIn;
    wrp->data = 0;
    wrp->ok = true;
    wrp->msg = "";

  } catch (const MacMidiError &err) {
    wrp->ptr = 0;
    wrp->data = 0;
    wrp->ok = false;
    wrp->msg = err.what();
  }

  return wrp;
}

void macmidi_in_free(MacMidiInPtr device) {
  if (device->data)
    delete (CallbackProxyUserData *)device->data;
  delete (MacMidiIn *)device->ptr;
  delete device;
}

enum MacMidiApi macmidi_in_get_current_api(MacMidiPtr device) {
  try {
    return (MacMidiApi)((MacMidiIn *)device->ptr)->getCurrentApi();

  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();

    return MACMIDI_API_UNSPECIFIED;
  }
}

static void callback_proxy(double timeStamp,
                           std::vector<unsigned char> *message,
                           void *userData) {
  CallbackProxyUserData *data =
      reinterpret_cast<CallbackProxyUserData *>(userData);
  data->c_callback(timeStamp, message->data(), message->size(),
                   data->user_data);
}

void macmidi_in_set_callback(MacMidiInPtr device, MacMidiCCallback callback,
                             void *userData) {
  device->data = (void *)new CallbackProxyUserData(callback, userData);
  try {
    ((MacMidiIn *)device->ptr)->setCallback(callback_proxy, device->data);
  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
    delete (CallbackProxyUserData *)device->data;
    device->data = 0;
  }
}

void macmidi_set_notification_callback(MIDINotifyProc callback) {
  createCoreMidiClientSingleton("MacMidiClient", callback);
}

void macmidi_in_cancel_callback(MacMidiInPtr device) {
  try {
    ((MacMidiIn *)device->ptr)->cancelCallback();
    delete (CallbackProxyUserData *)device->data;
    device->data = 0;
  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
  }
}

void macmidi_in_ignore_types(MacMidiInPtr device, bool midiSysex, bool midiTime,
                             bool midiSense) {
  ((MacMidiIn *)device->ptr)->ignoreTypes(midiSysex, midiTime, midiSense);
}

double macmidi_in_get_message(MacMidiInPtr device, unsigned char *message,
                              size_t *size) {
  try {
    // FIXME: use allocator to achieve efficient buffering
    std::vector<unsigned char> v;
    double ret = ((MacMidiIn *)device->ptr)->getMessage(&v);

    if (v.size() > 0 && v.size() <= *size) {
      memcpy(message, v.data(), (int)v.size());
    }

    *size = v.size();
    return ret;
  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
    return -1;
  } catch (...) {
    device->ok = false;
    device->msg = "Unknown error";
    return -1;
  }
}

/* MacMidiOut API */
MacMidiOutPtr macmidi_out_create_default() {
  MacMidiWrapper *wrp = new MacMidiWrapper;

  try {
    MacMidiOut *rOut = new MacMidiOut();

    wrp->ptr = (void *)rOut;
    wrp->data = 0;
    wrp->ok = true;
    wrp->msg = "";

  } catch (const MacMidiError &err) {
    wrp->ptr = 0;
    wrp->data = 0;
    wrp->ok = false;
    wrp->msg = err.what();
  }

  return wrp;
}

MacMidiOutPtr macmidi_out_create(enum MacMidiApi api, const char *clientName) {
  MacMidiWrapper *wrp = new MacMidiWrapper;
  std::string name = clientName;

  try {
    MacMidiOut *rOut = new MacMidiOut((MacMidi::Api)api, name);

    wrp->ptr = (void *)rOut;
    wrp->data = 0;
    wrp->ok = true;
    wrp->msg = "";

  } catch (const MacMidiError &err) {
    wrp->ptr = 0;
    wrp->data = 0;
    wrp->ok = false;
    wrp->msg = err.what();
  }

  return wrp;
}

void macmidi_out_free(MacMidiOutPtr device) {
  delete (MacMidiOut *)device->ptr;
  delete device;
}

enum MacMidiApi macmidi_out_get_current_api(MacMidiPtr device) {
  try {
    return (MacMidiApi)((MacMidiOut *)device->ptr)->getCurrentApi();

  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();

    return MACMIDI_API_UNSPECIFIED;
  }
}

int macmidi_out_send_message(MacMidiOutPtr device, const unsigned char *message,
                             int length) {
  try {
    ((MacMidiOut *)device->ptr)->sendMessage(message, length);
    return 0;
  } catch (const MacMidiError &err) {
    device->ok = false;
    device->msg = err.what();
    return -1;
  } catch (...) {
    device->ok = false;
    device->msg = "Unknown error";
    return -1;
  }
}
