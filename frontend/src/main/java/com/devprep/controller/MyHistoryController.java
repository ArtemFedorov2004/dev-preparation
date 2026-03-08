package com.devprep.controller;

import com.devprep.client.ApiClient;
import com.devprep.client.dto.ViewHistoryDto;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.util.List;

@Controller
@RequestMapping("/me/history")
@RequiredArgsConstructor
public class MyHistoryController {

    private static final String PAGE_TITLE_HISTORY = "История просмотров — DevPrep";

    private final ApiClient apiClient;

    @GetMapping
    public String history(Model model) {
        List<ViewHistoryDto> history = apiClient.getMyHistory();
        model.addAttribute("history", history);
        model.addAttribute("pageTitle", PAGE_TITLE_HISTORY);
        return "me/history";
    }
}
