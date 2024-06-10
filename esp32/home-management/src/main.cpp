#include <Arduino.h>
#include <ArduinoJson.h>
#include <EEPROM.h>
#include <PubSubClient.h>
#include <RCSwitch.h>
#include <StreamUtils.h>
#include <WiFi.h>

#include "time.h"

const int RF_READ_PIN = 13;
const int RF_WRITE_PIN = 14;
const int RF_CS_PIN = 12;
const int BLINK_INTERVAL = 500;
const int DEFAULT_FREQUENCY = 433;

const unsigned long BAUD_RATE = 115200;
const unsigned long EEPROM_SIZE = 512;
const unsigned long MQTT_BUFFER_SIZE = 1024;
const unsigned long DATE_BUFFER_SIZE = 80;
const long DEFAULT_RECEIVE_TOLERANCE = 60;

const String ACTION = "action";
const String ACTION_RESTART = "RESTART";
const String ACTION_RESET = "RESET";
const String ACTION_SET_CONFIG = "SET_CONFIG";
const String ACTION_GET_CONFIG = "GET_CONFIG";
const String ACTION_GET_INFO = "GET_INFO";
const String ACTION_SEND_RF = "SEND_RF";
const String ACTION_STATUS = "status";
const String STATUS_OK = "OK";
const String STATUS_ERROR = "ERROR";

const String CONFIG = "config";
const String FREQUENCY = "frequency";
const String TYPE = "type";
const String DEVICE = "device";
const String DATE = "date";
const String VERSION = "version";
const String CODE = "code";
const String BITS = "bits";
const String PROTOCOL = "protocol";
const String RECEIVE_TOLERANCE = "receiveTolerance";
const String UPTIME = "uptime";
const String WIFI = "wifi";
const String RSSI = "rssi";
const String IP = "ip";
const String SKETCH_MD5 = "sketchMD5";
const String SKETCH_SIZE = "sketchSize";
const String SDK_VERSION = "sdkVersion";

const String DEVICE_NAME = "esp32";
const String VERSION_VALUE = "1.0.0";
const String TYPE_RECEIVED_RF = "RECEIVED_RF";
const String TYPE_SEND_RF = "SEND_RF";
const String TYPE_MQTT_CONNECTED = "MQTT_CONNECTED";
const String TYPE_RESTART = "RESTART";
const String TYPE_RESET = "RESET";
const String TYPE_SET_CONFIG = "SET_CONFIG";
const String TYPE_GET_CONFIG = "GET_CONFIG";
const String TYPE_GET_INFO = "GET_INFO";
const String TYPE_ACTION_UNRECOGNIZED = "ACTION_UNRECOGNIZED";

const char* WIFI_NAME = "<REPLACE_ME>";
const char* WIFI_PASSWORD = "<REPLACE_ME>";

const char* MQTT_HOST = "<REPLACE_ME>";
const uint16_t MQTT_PORT = 1883;
const char* MQTT_USER = "esp32";
const char* MQTT_PASSWORD = "<REPLACE_ME>";

const char* TOPIC_BROADCAST = "home-management/broadcast";
const char* TOPIC_EVENTS = "home-management/events";
const char* TOPIC_DEVICES = "home-management/devices";
const char* TOPIC_DEVICES_SLASH = "home-management/devices/";

const char* NTP_SERVER = "br.pool.ntp.org";
const long GMT_OFFSET = -10800;  // GMT -3
const int DAYLIGHT_OFFSET = 0;
const int MQTT_QOS = 1;

unsigned long previousBlink = 0;

WiFiClient wifiClient;
PubSubClient mqttClient(wifiClient);

JsonDocument config;
RCSwitch mySwitch = RCSwitch();

void logInfo(String message) {
    Serial.println("[" + String(millis()) + "] INFO - " + message);
}

void logWarn(String message) {
    Serial.println("[" + String(millis()) + "] WARN - " + message);
}

void logError(String message) {
    Serial.println("[" + String(millis()) + "] ERROR - " + message);
}

void saveConfig() {
    EepromStream eepromStream(0, EEPROM_SIZE);
    serializeJson(config, eepromStream);
    eepromStream.flush();
    EEPROM.commit();
}

void loadConfig() {
    bool updateConfig = false;
    EepromStream eepromStream(0, EEPROM_SIZE);
    deserializeJson(config, eepromStream);

    if (!config.containsKey(RECEIVE_TOLERANCE)) {
        config[RECEIVE_TOLERANCE] = DEFAULT_RECEIVE_TOLERANCE;
    }

    if (updateConfig) {
        saveConfig();
    }
}

long getReceiveTolerance() {
    return config[RECEIVE_TOLERANCE];
}

String getCurrentFormattedDate() {
    struct tm timeinfo;

    if (getLocalTime(&timeinfo)) {
        char buffer[DATE_BUFFER_SIZE];
        strftime(buffer, sizeof(buffer), "%FT%T-03:00", &timeinfo);
        return String(buffer);
    }

    return "";
}

JsonDocument addExtraDataResponseJson(JsonDocument responseJson) {
    String currentFormattedDate = getCurrentFormattedDate();
    if (currentFormattedDate != "") {
        responseJson[DATE] = currentFormattedDate;
    }

    responseJson[DEVICE] = DEVICE_NAME;
    responseJson[VERSION] = VERSION_VALUE;
    return responseJson;
}

bool checkWifiConnected() {
    bool connected = false;

    for (int i = 0; i < 10; i++) {
        connected = WiFi.isConnected();
        if (connected) {
            break;
        } else {
            delay(1000);
        }
    }

    return connected;
}

String getDeviceTopic() {
    return TOPIC_DEVICES_SLASH + DEVICE_NAME;
}

void sendMqttConnectedMessage() {
    JsonDocument eventJson;
    eventJson[TYPE] = TYPE_MQTT_CONNECTED;

    eventJson[WIFI][RSSI] = WiFi.RSSI();
    eventJson[WIFI][IP] = WiFi.localIP();

    eventJson[SKETCH_MD5] = ESP.getSketchMD5();
    eventJson[SKETCH_SIZE] = ESP.getSketchSize();
    eventJson[SDK_VERSION] = ESP.getSdkVersion();
    eventJson[UPTIME] = millis();

    eventJson = addExtraDataResponseJson(eventJson);

    char buffer[MQTT_BUFFER_SIZE];
    size_t n = serializeJson(eventJson, buffer);

    mqttClient.publish(TOPIC_EVENTS, buffer, n);
}

bool mqttConnect() {
    bool connected = false;
    if (WiFi.isConnected() && !mqttClient.connected()) {
        logInfo("MQTT connecting...");

        mqttClient.setServer(MQTT_HOST, MQTT_PORT);
        mqttClient.setBufferSize(MQTT_BUFFER_SIZE);

        if (mqttClient.connect(DEVICE_NAME.c_str(), MQTT_USER, MQTT_PASSWORD, 0, MQTT_QOS, false, 0, false)) {
            logInfo("MQTT connected");

            String topic = getDeviceTopic();
            mqttClient.subscribe(topic.c_str(), MQTT_QOS);
            mqttClient.subscribe(TOPIC_BROADCAST, MQTT_QOS);
            mqttClient.subscribe(TOPIC_DEVICES, MQTT_QOS);
            connected = true;

            sendMqttConnectedMessage();
        }
    }

    return connected;
}

void syncNtp() {
    configTime(GMT_OFFSET, DAYLIGHT_OFFSET, NTP_SERVER);

    while (getCurrentFormattedDate() == "") {
        logInfo("Waiting for NTP sync...");
        delay(1000);
    }
}

void wifiGotIp(WiFiEvent_t event) {
    if (event == ARDUINO_EVENT_WIFI_STA_GOT_IP) {
        syncNtp();
        mqttConnect();
    }
}

void setupWifi() {
    WiFi.disconnect(true);

    WiFi.config(INADDR_NONE, INADDR_NONE, INADDR_NONE);
    WiFi.setHostname(DEVICE_NAME.c_str());
    WiFi.mode(WIFI_STA);
    WiFi.onEvent(wifiGotIp);
    WiFi.begin(WIFI_NAME, WIFI_PASSWORD);
}

void restartBackground(void* arg) {
    delay(1000);
    ESP.restart();
}

void restart() {
    xTaskCreatePinnedToCore(restartBackground,
                            "restartBackground",
                            2048,
                            NULL,
                            10,
                            NULL,
                            tskNO_AFFINITY);
}

void sendMqttMessage(String topic, JsonDocument responseJson) {
    if (!responseJson.isNull()) {
        bool published = false;
        responseJson = addExtraDataResponseJson(responseJson);

        for (int i = 0; i < 10; i++) {
            char buffer[MQTT_BUFFER_SIZE];
            size_t n = serializeJson(responseJson, buffer);

            published = mqttClient.publish(topic.c_str(), buffer, n);
            if (published) {
                break;
            } else {
                logInfo("MQTT message not sended, retrying... " + String(i));
                delay(1000);
                mqttConnect();
            }
        }

        if (published) {
            logInfo("MQTT message sended");
        } else {
            logError("MQTT message not sended");
        }
    }
}

void handleActionGetInfo() {
    JsonDocument eventJson;
    eventJson[TYPE] = TYPE_GET_INFO;
    eventJson[ACTION_STATUS] = STATUS_OK;

    eventJson[WIFI][RSSI] = WiFi.RSSI();
    eventJson[WIFI][IP] = WiFi.localIP();

    eventJson[SKETCH_MD5] = ESP.getSketchMD5();
    eventJson[SKETCH_SIZE] = ESP.getSketchSize();
    eventJson[SDK_VERSION] = ESP.getSdkVersion();
    eventJson[UPTIME] = millis();

    sendMqttMessage(TOPIC_EVENTS, eventJson);
}

void handleActionGetConfig() {
    JsonDocument eventJson;
    eventJson[TYPE] = TYPE_GET_CONFIG;
    eventJson[ACTION_STATUS] = STATUS_OK;
    eventJson[CONFIG] = config;

    sendMqttMessage(TOPIC_EVENTS, eventJson);
}

void handleActionSetConfig(JsonDocument receivedJson) {
    config = receivedJson[CONFIG];
    saveConfig();
    restart();

    JsonDocument eventJson;
    eventJson[TYPE] = TYPE_SET_CONFIG;
    eventJson[ACTION_STATUS] = STATUS_OK;

    sendMqttMessage(TOPIC_EVENTS, eventJson);
}

void handleActionResetDevice() {
    for (int i = 0; i < EEPROM_SIZE; i++) {
        EEPROM.write(i, 0);
    }

    EEPROM.commit();
    restart();

    JsonDocument eventJson;
    eventJson[TYPE] = TYPE_RESET;
    eventJson[ACTION_STATUS] = STATUS_OK;

    sendMqttMessage(TOPIC_EVENTS, eventJson);
}

void handleActionRestartDevice() {
    restart();

    JsonDocument eventJson;
    eventJson[TYPE] = TYPE_RESTART;
    eventJson[ACTION_STATUS] = STATUS_OK;

    sendMqttMessage(TOPIC_EVENTS, eventJson);
}

void handleActionSendRF(JsonDocument receivedJson) {
    int protocol = receivedJson[PROTOCOL];
    unsigned long code = receivedJson[CODE];
    unsigned int length = receivedJson[BITS];

    mySwitch.setProtocol(protocol);
    mySwitch.send(code, length);

    JsonDocument eventJson;
    eventJson[TYPE] = TYPE_SEND_RF;
    eventJson[ACTION_STATUS] = STATUS_OK;

    sendMqttMessage(TOPIC_EVENTS, eventJson);
}

void handleActionUnrecognized(String action) {
    JsonDocument eventJson;
    eventJson[TYPE] = TYPE_ACTION_UNRECOGNIZED;
    eventJson[ACTION_STATUS] = STATUS_ERROR;
    eventJson[ACTION] = action;

    sendMqttMessage(TOPIC_EVENTS, eventJson);
}

void mqttCallback(char* topic, byte* message, unsigned int length) {
    String topicStr = String(topic);
    String deviceTopic = getDeviceTopic();

    if (topicStr == deviceTopic) {
        digitalWrite(LED_BUILTIN, HIGH);
        logInfo("MQTT message received from topic [" + topicStr + "]");

        JsonDocument receivedJson;
        deserializeJson(receivedJson, message, length);

        if (!receivedJson.isNull()) {
            String action = receivedJson[ACTION];
            if (ACTION_GET_INFO.equals(action)) {
                handleActionGetInfo();

            } else if (ACTION_GET_CONFIG.equals(action)) {
                handleActionGetConfig();

            } else if (ACTION_SET_CONFIG.equals(action)) {
                handleActionSetConfig(receivedJson);

            } else if (ACTION_RESTART.equals(action)) {
                handleActionRestartDevice();

            } else if (ACTION_RESET.equals(action)) {
                handleActionResetDevice();

            } else if (ACTION_SEND_RF.equals(action)) {
                handleActionSendRF(receivedJson);

            } else if (action != "null") {
                handleActionUnrecognized(action);
            }
        }

        logInfo("MQTT message consumed from topic [" + topicStr + "]");
        digitalWrite(LED_BUILTIN, LOW);
    }
}

void setup() {
    Serial.begin(115200);
    logInfo("ESP32 starting...");

    pinMode(LED_BUILTIN, OUTPUT);
    pinMode(RF_CS_PIN, OUTPUT);
    digitalWrite(RF_CS_PIN, HIGH);

    EEPROM.begin(EEPROM_SIZE);
    loadConfig();

    long receiveTolerance = getReceiveTolerance();
    mySwitch.enableReceive(RF_READ_PIN);
    mySwitch.enableTransmit(RF_WRITE_PIN);
    mySwitch.setReceiveTolerance(receiveTolerance);

    mqttClient.setCallback(mqttCallback);

    setupWifi();

    while (!mqttClient.connected()) {
        logInfo("Waiting for MQTT connect...");
        delay(1000);
    }

    logInfo("ESP32 started");
}

void loop() {
    unsigned long currentMillis = millis();

    if (currentMillis - previousBlink >= BLINK_INTERVAL) {
        previousBlink = currentMillis;

        uint8_t state = !digitalRead(LED_BUILTIN);
        digitalWrite(LED_BUILTIN, state);
    }

    boolean connected = mqttClient.loop();
    if (!connected) {
        mqttConnect();
    }

    if (mySwitch.available()) {
        digitalWrite(LED_BUILTIN, HIGH);

        long code = mySwitch.getReceivedValue();
        int bits = mySwitch.getReceivedBitlength();
        int protocol = mySwitch.getReceivedProtocol();
        long receiveTolerance = getReceiveTolerance();

        mySwitch.resetAvailable();

        JsonDocument eventJson;
        eventJson[TYPE] = TYPE_RECEIVED_RF;
        eventJson[CODE] = code;
        eventJson[BITS] = bits;
        eventJson[PROTOCOL] = protocol;
        eventJson[RECEIVE_TOLERANCE] = receiveTolerance;
        eventJson[FREQUENCY] = DEFAULT_FREQUENCY;
        sendMqttMessage(TOPIC_EVENTS, eventJson);

        digitalWrite(LED_BUILTIN, LOW);
    }

    if (Serial.available() > 0) {
        String input = Serial.readStringUntil('\n');
        logInfo("Received from serial: " + input);
    }
}
