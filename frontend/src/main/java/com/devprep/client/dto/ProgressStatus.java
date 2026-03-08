package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonValue;

public enum ProgressStatus {
    LEARNED("learned"),
    NEED_REVIEW("need_review"),
    DONT_KNOW("dont_know");

    private final String value;

    ProgressStatus(String value) {
        this.value = value;
    }

    @JsonValue
    public String getValue() {
        return value;
    }

    @JsonCreator
    public static ProgressStatus fromValue(String value) {
        for (ProgressStatus status : ProgressStatus.values()) {
            if (status.value.equalsIgnoreCase(value)) {
                return status;
            }
        }
        throw new IllegalArgumentException("Unknown progress status: " + value);
    }

    public String getDisplayName() {
        return switch (this) {
            case LEARNED -> "Изучено";
            case NEED_REVIEW -> "Повторить";
            case DONT_KNOW -> "Не знаю";
        };
    }

    public String getBadgeClass() {
        return switch (this) {
            case LEARNED -> "bg-green-100 text-green-700";
            case NEED_REVIEW -> "bg-yellow-100 text-yellow-700";
            case DONT_KNOW -> "bg-red-100 text-red-700";
        };
    }
}