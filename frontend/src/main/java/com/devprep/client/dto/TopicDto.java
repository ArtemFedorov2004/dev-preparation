package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class TopicDto {

    private int id;

    private String slug;

    private String name;

    private String description;

    private String icon;

    @JsonProperty("sort_order")
    private int sortOrder;
}