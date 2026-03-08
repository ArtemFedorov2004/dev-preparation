package com.devprep.client;

import com.devprep.client.dto.*;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Primary;
import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Component;

import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.*;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;

@Slf4j
@Component
@Profile({"mock", "test"})
@Primary
public class MockApiClient implements ApiClient {

    private final List<TopicDto> mockTopics = new ArrayList<>();

    private final List<TagDto> mockTags = new ArrayList<>();

    private final Map<String, QuestionDetailDto> questionsBySlug = new LinkedHashMap<>();

    private final Map<String, List<QuestionListItemDto>> questionsByTopic = new HashMap<>();

    private final Map<String, ProgressStatus> progressStore = new ConcurrentHashMap<>();

    private final Map<String, Instant> progressUpdatedAt = new ConcurrentHashMap<>();

    private final Map<String, Boolean> bookmarkStore = new ConcurrentHashMap<>();

    private final Map<String, Instant> bookmarkAddedAt = new ConcurrentHashMap<>();

    private final LinkedList<String> viewHistory = new LinkedList<>();

    private final Map<String, Instant> lastViewedAt = new ConcurrentHashMap<>();

    private final AtomicInteger idSequence = new AtomicInteger(1);

    public MockApiClient() {
        createTags();
        createTopics();
        createQuestions();
        seed();
        log.info("MockApiClient: {} topics, {} tags, {} questions",
                mockTopics.size(), mockTags.size(), questionsBySlug.size());
    }

    private void seed() {
        List<String> slugs = new ArrayList<>(questionsBySlug.keySet());
        if (slugs.isEmpty()) return;

        setProgressInternal(slugs.get(0), ProgressStatus.LEARNED, Instant.now().minus(5, ChronoUnit.DAYS));
        if (slugs.size() > 1)
            setProgressInternal(slugs.get(1), ProgressStatus.LEARNED, Instant.now().minus(3, ChronoUnit.DAYS));
        if (slugs.size() > 2)
            setProgressInternal(slugs.get(2), ProgressStatus.NEED_REVIEW, Instant.now().minus(1, ChronoUnit.DAYS));
        if (slugs.size() > 3)
            setProgressInternal(slugs.get(3), ProgressStatus.DONT_KNOW, Instant.now().minus(2, ChronoUnit.DAYS));

        addBookmarkInternal(slugs.get(1), Instant.now().minus(3, ChronoUnit.DAYS));
        if (slugs.size() > 4) addBookmarkInternal(slugs.get(4), Instant.now().minus(1, ChronoUnit.DAYS));

        for (int i = Math.min(slugs.size(), 5) - 1; i >= 0; i--) {
            addViewInternal(slugs.get(i), Instant.now().minus(i, ChronoUnit.HOURS));
        }
    }

    private void setProgressInternal(String slug, ProgressStatus status, Instant at) {
        progressStore.put(slug, status);
        progressUpdatedAt.put(slug, at);
    }

    private void addBookmarkInternal(String slug, Instant at) {
        bookmarkStore.put(slug, true);
        bookmarkAddedAt.put(slug, at);
    }

    private void addViewInternal(String slug, Instant at) {
        viewHistory.remove(slug);
        viewHistory.addFirst(slug);
        lastViewedAt.put(slug, at);
    }

    @Override
    public List<TopicDto> getTopics() {
        return new ArrayList<>(mockTopics);
    }

    @Override
    public Optional<TopicWithQuestionsDto> getTopicBySlug(String slug, Level level) {
        return mockTopics.stream()
                .filter(t -> t.slug().equals(slug))
                .findFirst()
                .map(topic -> {
                    TopicWithQuestionsDto dto = new TopicWithQuestionsDto();
                    dto.setId(topic.id());
                    dto.setSlug(topic.slug());
                    dto.setName(topic.name());
                    dto.setDescription(topic.description());
                    dto.setIcon(topic.icon());
                    dto.setSortOrder(topic.sortOrder());
                    dto.setQuestions(
                            questionsByTopic.getOrDefault(slug, Collections.emptyList()).stream()
                                    .filter(q -> level == null || level.equals(q.level()))
                                    .collect(Collectors.toList()));
                    return dto;
                });
    }

    @Override
    public PaginatedQuestionsDto getQuestions(String topic, String tag, Level level, int page, int limit) {
        List<QuestionListItemDto> filtered = questionsBySlug.values().stream()
                .map(this::toListItem)
                .collect(Collectors.toList());

        if (topic != null && !topic.isBlank())
            filtered = filtered.stream()
                    .filter(q -> q.topic() != null && topic.equals(q.topic().slug()))
                    .collect(Collectors.toList());

        if (tag != null && !tag.isBlank())
            filtered = filtered.stream()
                    .filter(q -> q.tags() != null &&
                            q.tags().stream().anyMatch(t -> tag.equals(t.slug())))
                    .collect(Collectors.toList());

        if (level != null)
            filtered = filtered.stream()
                    .filter(q -> level.equals(q.level()))
                    .collect(Collectors.toList());

        int total = filtered.size();
        int fromIndex = Math.min((page - 1) * limit, total);
        int toIndex = Math.min(fromIndex + limit, total);

        PaginatedQuestionsDto result = new PaginatedQuestionsDto();
        result.setData(fromIndex < toIndex ? filtered.subList(fromIndex, toIndex) : Collections.emptyList());
        PaginatedQuestionsDto.PaginationDto pagination = new PaginatedQuestionsDto.PaginationDto();
        pagination.setPage(page);
        pagination.setLimit(limit);
        pagination.setTotal(total);
        result.setPagination(pagination);
        return result;
    }

    @Override
    public Optional<QuestionDetailDto> getQuestionBySlug(String slug) {
        return Optional.ofNullable(questionsBySlug.get(slug));
    }

    @Override
    public List<TagDto> getTags() {
        return new ArrayList<>(mockTags);
    }

    @Override
    public void updateProgress(String slug, ProgressStatus status) {
        log.debug("Mock: updateProgress slug={} status={}", slug, status);
        setProgressInternal(slug, status, Instant.now());
    }

    @Override
    public List<UserProgressDto> getMyProgress() {
        return progressStore.entrySet().stream()
                .map(e -> {
                    QuestionDetailDto q = questionsBySlug.get(e.getKey());
                    if (q == null) return null;
                    UserProgressDto dto = new UserProgressDto();
                    dto.setSlug(q.slug());
                    dto.setTitle(q.title());
                    dto.setLevel(q.level());
                    dto.setTopic(q.topic());
                    dto.setStatus(e.getValue());
                    dto.setUpdatedAt(progressUpdatedAt.getOrDefault(e.getKey(), Instant.now()));
                    return dto;
                })
                .filter(Objects::nonNull)
                .sorted(Comparator.comparing(UserProgressDto::getUpdatedAt).reversed())
                .collect(Collectors.toList());
    }

    @Override
    public List<TopicProgressDto> getMyProgressByTopic() {
        Map<String, TopicProgressDto> result = new LinkedHashMap<>();

        for (TopicDto topic : mockTopics) {
            List<QuestionListItemDto> questions =
                    questionsByTopic.getOrDefault(topic.slug(), Collections.emptyList());
            if (questions.isEmpty()) continue;

            int learned = 0, needReview = 0, dontKnow = 0;
            for (QuestionListItemDto q : questions) {
                ProgressStatus s = progressStore.get(q.slug());
                if (s == null) continue;
                switch (s) {
                    case LEARNED -> learned++;
                    case NEED_REVIEW -> needReview++;
                    case DONT_KNOW -> dontKnow++;
                }
            }

            TopicProgressDto tp = new TopicProgressDto(
                    topic,
                    questions.size(),
                    learned,
                    needReview,
                    dontKnow
            );

            result.put(topic.slug(), tp);
        }

        return new ArrayList<>(result.values());
    }

    @Override
    public Optional<ProgressStatus> getQuestionProgress(String slug) {
        return Optional.ofNullable(progressStore.get(slug));
    }

    @Override
    public boolean toggleBookmark(String slug) {
        boolean current = bookmarkStore.getOrDefault(slug, false);
        boolean next = !current;
        if (next) {
            bookmarkStore.put(slug, true);
            bookmarkAddedAt.put(slug, Instant.now());
        } else {
            bookmarkStore.remove(slug);
            bookmarkAddedAt.remove(slug);
        }
        log.debug("Mock: bookmark {} slug={}", next ? "ADDED" : "REMOVED", slug);
        return next;
    }

    @Override
    public List<BookmarkDto> getMyBookmarks() {
        return bookmarkStore.entrySet().stream()
                .filter(Map.Entry::getValue)
                .map(e -> {
                    QuestionDetailDto q = questionsBySlug.get(e.getKey());
                    if (q == null) return null;
                    BookmarkDto dto = new BookmarkDto();
                    dto.setSlug(q.slug());
                    dto.setTitle(q.title());
                    dto.setLevel(q.level());
                    dto.setTopic(q.topic());
                    dto.setBookmarkedAt(bookmarkAddedAt.getOrDefault(e.getKey(), Instant.now()));
                    return dto;
                })
                .filter(Objects::nonNull)
                .sorted(Comparator.comparing(BookmarkDto::getBookmarkedAt).reversed())
                .collect(Collectors.toList());
    }

    @Override
    public boolean isBookmarked(String slug) {
        return bookmarkStore.getOrDefault(slug, false);
    }

    @Override
    public void recordView(String slug) {
        log.debug("Mock: recordView slug={}", slug);
        addViewInternal(slug, Instant.now());
        while (viewHistory.size() > 50) {
            String removed = viewHistory.removeLast();
            lastViewedAt.remove(removed);
        }
    }

    @Override
    public List<ViewHistoryDto> getMyHistory() {
        return viewHistory.stream()
                .map(slug -> {
                    QuestionDetailDto q = questionsBySlug.get(slug);
                    if (q == null) return null;
                    ViewHistoryDto dto = new ViewHistoryDto();
                    dto.setSlug(q.slug());
                    dto.setTitle(q.title());
                    dto.setLevel(q.level());
                    dto.setTopic(q.topic());
                    dto.setViewedAt(lastViewedAt.getOrDefault(slug, Instant.now()));
                    return dto;
                })
                .filter(Objects::nonNull)
                .collect(Collectors.toList());
    }

    private void createTags() {
        mockTags.add(tag("java", "Java"));
        mockTags.add(tag("python", "Python"));
        mockTags.add(tag("javascript", "JavaScript"));
        mockTags.add(tag("go", "Go"));
        mockTags.add(tag("rust", "Rust"));
        mockTags.add(tag("oop", "ООП"));
        mockTags.add(tag("functional", "Функциональное программирование"));
        mockTags.add(tag("concurrency", "Параллелизм"));
        mockTags.add(tag("memory-management", "Управление памятью"));
        mockTags.add(tag("algorithms", "Алгоритмы"));
        mockTags.add(tag("data-structures", "Структуры данных"));
        mockTags.add(tag("design-patterns", "Шаблоны проектирования"));
        mockTags.add(tag("sql", "SQL"));
        mockTags.add(tag("nosql", "NoSQL"));
        mockTags.add(tag("testing", "Тестирование"));
        mockTags.add(tag("microservices", "Микросервисы"));
        mockTags.add(tag("docker", "Docker"));
        mockTags.add(tag("kubernetes", "Kubernetes"));
        mockTags.add(tag("security", "Безопасность"));
    }

    private void createTopics() {
        mockTopics.add(topic("oop", "ООП", "Основные концепции ООП: классы, объекты, наследование, полиморфизм, инкапсуляция", "🔷", 1));
        mockTopics.add(topic("java-core", "Основы Java", "Базовые концепции Java, синтаксис и стандартные библиотеки", "☕", 2));
        mockTopics.add(topic("databases", "Базы данных", "SQL, NoSQL, индексы, транзакции и проектирование баз данных", "🗄️", 3));
        mockTopics.add(topic("algorithms", "Алгоритмы", "Распространённые алгоритмы, анализ сложности и реализация структур данных", "⚙️", 4));
        mockTopics.add(topic("spring-boot", "Spring Boot", "Фреймворк Spring, внедрение зависимостей, веб-разработка и микросервисы", "🍃", 5));
        mockTopics.add(topic("system-design", "Проектирование систем", "Проектирование масштабируемых систем, архитектурные шаблоны и компромиссы", "🏗️", 6));
        mockTopics.add(topic("devops", "DevOps и облака", "CI/CD, контейнеризация, облачные платформы и инфраструктура как код", "🚀", 7));
        mockTopics.add(topic("javascript", "JavaScript", "Современный JavaScript, ES6+, браузерные API и Node.js", "📜", 8));
    }

    private void createQuestions() {
        q("what-is-polymorphism",
                "What is polymorphism and what are its types in OOP?",
                """
                        # Polymorphism in OOP
                        
                        **Polymorphism** allows objects of different classes to be treated as objects of a common superclass.
                        
                        ## Types
                        
                        ### 1. Compile-time (Method Overloading)
                        Same method name, different parameters.
                        
                        ### 2. Runtime (Method Overriding)
                        Subclass overrides a method from its superclass.
                        
                        ```java
                        abstract class Animal { public abstract void makeSound(); }
                        class Dog extends Animal { public void makeSound() { System.out.println("Woof!"); } }
                        class Cat extends Animal { public void makeSound() { System.out.println("Meow!"); } }
                        
                        Animal pet = new Dog();
                        pet.makeSound(); // Woof!
                        ```
                        """,
                Level.JUNIOR, "oop", List.of("oop", "java", "design-patterns"));

        q("inheritance-vs-composition",
                "Inheritance vs Composition: Which should you prefer and why?",
                """
                        # Inheritance vs Composition
                        
                        Modern best practices favor **composition over inheritance**.
                        
                        ## Inheritance — "is-a"
                        Tight coupling, fragile base class problem, deep hierarchies become hard to maintain.
                        
                        ## Composition — "has-a"
                        ```java
                        class Engine { public void start() { ... } }
                        class Car {
                            private final Engine engine = new Engine();
                            public void start() { engine.start(); }
                        }
                        ```
                        Loose coupling, easier to test, more flexible.
                        """,
                Level.MIDDLE, "oop", List.of("oop", "design-patterns", "java"));

        q("solid-principles-explained",
                "Can you explain the SOLID principles with examples?",
                """
                        # SOLID Principles
                        
                        - **S** — Single Responsibility: one reason to change.
                        - **O** — Open/Closed: open for extension, closed for modification.
                        - **L** — Liskov Substitution: subtypes must be substitutable for their base types.
                        - **I** — Interface Segregation: no client should depend on methods it doesn't use.
                        - **D** — Dependency Inversion: depend on abstractions, not concretions.
                        """,
                Level.SENIOR, "oop", List.of("oop", "design-patterns", "java"));

        q("java-memory-management",
                "How does memory management work in Java?",
                """
                        # Java Memory Management
                        
                        ## JVM Memory Structure
                        - **Heap** — all objects; Young Gen (Eden + Survivor) + Old Gen.
                        - **Stack** — per-thread; primitives, references, call frames.
                        - **Metaspace** — class metadata, static variables.
                        
                        ## GC Algorithms
                        - **G1 GC** (default since Java 9) — region-based, predictable pauses.
                        - **ZGC** (Java 11+) — sub-10ms pauses.
                        
                        ```java
                        try (var fis = new FileInputStream("file.txt")) { /* auto-closed */ }
                        ```
                        """,
                Level.MIDDLE, "java-core", List.of("java", "memory-management"));

        q("java-concurrency-basics",
                "What are the different ways to create threads in Java?",
                """
                        # Thread Creation in Java
                        
                        ```java
                        // 1. Runnable (preferred for simple tasks)
                        new Thread(() -> doWork()).start();
                        
                        // 2. ExecutorService (production)
                        var exec = Executors.newFixedThreadPool(5);
                        exec.submit(() -> doWork());
                        
                        // 3. CompletableFuture (async pipelines)
                        CompletableFuture.supplyAsync(this::fetch)
                            .thenApply(this::process);
                        
                        // 4. Virtual Threads — Java 21+
                        try (var exec = Executors.newVirtualThreadPerTaskExecutor()) {
                            exec.submit(() -> doWork());
                        }
                        ```
                        """,
                Level.MIDDLE, "java-core", List.of("java", "concurrency"));

        q("sql-vs-nosql",
                "SQL vs NoSQL: When to use which?",
                """
                        # SQL vs NoSQL
                        
                        | | SQL | NoSQL |
                        |---|---|---|
                        | Schema | Fixed | Flexible |
                        | Transactions | ACID | Eventual consistency |
                        | Scale | Vertical | Horizontal |
                        | Best for | Finance, orders | Social feeds, caching |
                        
                        **Polyglot persistence** — combine both in one system.
                        """,
                Level.MIDDLE, "databases", List.of("sql", "nosql"));

        q("rest-api-design-principles",
                "What are the key principles of REST API design?",
                """
                        # REST API Design
                        
                        1. **Nouns, not verbs**: `GET /users/123`, not `GET /getUser?id=123`
                        2. **HTTP methods**: GET read, POST create, PUT replace, PATCH update, DELETE remove
                        3. **Status codes**: 200, 201, 204, 400, 401, 404, 500
                        4. **Pagination**: `GET /questions?page=1&limit=20`
                        5. **Versioning**: `/api/v1/`
                        """,
                Level.MIDDLE, "system-design", List.of("microservices", "java"));
    }

    private TagDto tag(String slug, String name) {
        TagDto t = new TagDto(
                idSequence.getAndIncrement(),
                slug,
                name
        );
        return t;
    }

    private TopicDto topic(String slug, String name, String description, String icon, int sortOrder) {
        TopicDto t = new TopicDto(
                idSequence.getAndIncrement(),
                slug,
                name,
                description,
                icon,
                sortOrder
        );
        return t;
    }

    private TagDto findTag(String slug) {
        return mockTags.stream().filter(t -> t.slug().equals(slug)).findFirst()
                .orElseThrow(() -> new IllegalArgumentException("Tag not found: " + slug));
    }

    private TopicDto findTopic(String slug) {
        return mockTopics.stream().filter(t -> t.slug().equals(slug)).findFirst()
                .orElseThrow(() -> new IllegalArgumentException("Topic not found: " + slug));
    }

    private void q(String slug, String title, String answer, Level level,
                   String topicSlug, List<String> tagSlugs) {
        TopicDto topic = findTopic(topicSlug);
        List<TagDto> tags = tagSlugs.stream().map(this::findTag).collect(Collectors.toList());

        QuestionDetailDto detail = new QuestionDetailDto(
                idSequence.getAndIncrement(),
                slug,
                title,
                answer,
                level,
                topic,
                tags
        );

        questionsBySlug.put(slug, detail);
        questionsByTopic.computeIfAbsent(topicSlug, k -> new ArrayList<>()).add(toListItem(detail));
    }

    private QuestionListItemDto toListItem(QuestionDetailDto d) {
        QuestionListItemDto item = new QuestionListItemDto(
                d.id(),
                d.slug(),
                d.title(),
                d.level(),
                d.topic(),
                d.tags()
        );

        return item;
    }
}