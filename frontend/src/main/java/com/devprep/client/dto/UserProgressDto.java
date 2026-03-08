package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.annotation.JsonProperty;
import lombok.Data;

import java.time.Instant;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class UserProgressDto {

    @JsonProperty("question_id")
    private int questionId;

    private String slug;

    private String title;

    private Level level;

    private TopicDto topic;

    private ProgressStatus status;

    @JsonProperty("updated_at")
    private Instant updatedAt;
}