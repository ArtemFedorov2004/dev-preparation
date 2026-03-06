package com.devprep.client.dto;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import lombok.Data;

import java.util.List;

@Data
@JsonIgnoreProperties(ignoreUnknown = true)
public class PaginatedQuestionsDto {

    private List<QuestionListItemDto> data;

    private PaginationDto pagination;

    @Data
    @JsonIgnoreProperties(ignoreUnknown = true)
    public static class PaginationDto {

        private int page;

        private int limit;

        private int total;
    }
}
