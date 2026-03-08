package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;

@JsonIgnoreProperties(ignoreUnknown = true)
public record TopicDto(
        int id, String slug, String name, String description, String icon,
        @JsonProperty("sort_order")
        int sortOrder
) {
}