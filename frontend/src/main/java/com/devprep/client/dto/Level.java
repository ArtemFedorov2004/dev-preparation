package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonValue;

public enum Level {
    JUNIOR("junior"),
    MIDDLE("middle"),
    SENIOR("senior");

    private final String value;

    Level(String value) {
        this.value = value;
    }

    @JsonValue
    public String getValue() {
        return value;
    }

    @JsonCreator
    public static Level fromValue(String value) {
        for (Level level : Level.values()) {
            if (level.value.equalsIgnoreCase(value)) {
                return level;
            }
        }
        throw new IllegalArgumentException("Unknown level: " + value);
    }

    public String getDisplayName() {
        return value.substring(0, 1).toUpperCase() + value.substring(1);
    }
}
