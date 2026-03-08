package com.devprep.controller;

import com.devprep.client.ApiClient;
import com.devprep.client.dto.ProgressStatus;
import com.devprep.client.dto.TopicProgressDto;
import com.devprep.client.dto.UpdateProgressRequest;
import com.devprep.client.dto.UserProgressDto;
import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;

import java.util.List;

@Controller
@RequestMapping("/me")
@RequiredArgsConstructor
public class MyProgressController {

    private static final String PAGE_TITLE_PROGRESS = "Мой прогресс — DevPrep";

    private final ApiClient apiClient;

    @GetMapping("/progress")
    public String progressDashboard(Model model) {
        List<TopicProgressDto> progressByTopic = apiClient.getMyProgressByTopic();
        List<UserProgressDto> progressList = apiClient.getMyProgress();

        model.addAttribute("progressByTopic", progressByTopic);
        model.addAttribute("progressList", progressList);
        model.addAttribute("pageTitle", PAGE_TITLE_PROGRESS);
        return "me/progress";
    }

    @PostMapping("/questions/{slug}/progress")
    public String updateProgress(
            @PathVariable String slug,
            UpdateProgressRequest request) {
        ProgressStatus status = ProgressStatus.fromValue(request.getStatus());
        apiClient.updateProgress(slug, status);
        return "redirect:/questions/%s".formatted(slug);
    }
}
