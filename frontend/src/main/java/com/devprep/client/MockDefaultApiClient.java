package com.devprep.client;

import com.devprep.client.dto.*;
import lombok.extern.slf4j.Slf4j;
import org.springframework.context.annotation.Primary;
import org.springframework.context.annotation.Profile;
import org.springframework.stereotype.Component;

import java.util.*;
import java.util.concurrent.atomic.AtomicInteger;
import java.util.stream.Collectors;

@Slf4j
@Component
@Profile({"mock", "test"})
@Primary
public class MockDefaultApiClient implements ApiClient {

    private final List<TopicDto> mockTopics = new ArrayList<>();
    private final List<TagDto> mockTags = new ArrayList<>();
    private final Map<Integer, QuestionDetailDto> mockQuestionsById = new LinkedHashMap<>();
    private final Map<String, List<QuestionListItemDto>> questionsByTopicSlug = new HashMap<>();

    private final AtomicInteger idSequence = new AtomicInteger(1);

    public MockDefaultApiClient() {
        initializeMockData();
        log.info("MockApiClient initialized with {} topics, {} tags, and {} questions",
                mockTopics.size(), mockTags.size(), mockQuestionsById.size());
    }

    private void initializeMockData() {
        createTags();
        createTopics();
        createQuestions();
    }

    private void createTags() {
        mockTags.add(createTag("java", "Java"));
        mockTags.add(createTag("python", "Python"));
        mockTags.add(createTag("javascript", "JavaScript"));
        mockTags.add(createTag("go", "Go"));
        mockTags.add(createTag("rust", "Rust"));
        mockTags.add(createTag("oop", "ООП"));
        mockTags.add(createTag("functional", "Функциональное программирование"));
        mockTags.add(createTag("concurrency", "Параллелизм"));
        mockTags.add(createTag("memory-management", "Управление памятью"));
        mockTags.add(createTag("algorithms", "Алгоритмы"));
        mockTags.add(createTag("data-structures", "Структуры данных"));
        mockTags.add(createTag("design-patterns", "Шаблоны проектирования"));
        mockTags.add(createTag("sql", "SQL"));
        mockTags.add(createTag("nosql", "NoSQL"));
        mockTags.add(createTag("testing", "Тестирование"));
        mockTags.add(createTag("microservices", "Микросервисы"));
        mockTags.add(createTag("docker", "Docker"));
        mockTags.add(createTag("kubernetes", "Kubernetes"));
        mockTags.add(createTag("security", "Безопасность"));
    }

    private void createTopics() {
        mockTopics.add(createTopic("oop", "ООП",
                "Основные концепции ООП: классы, объекты, наследование, полиморфизм, инкапсуляция", "🔷", 1));
        mockTopics.add(createTopic("java-core", "Основы Java",
                "Базовые концепции Java, синтаксис и стандартные библиотеки", "☕", 2));
        mockTopics.add(createTopic("databases", "Базы данных",
                "SQL, NoSQL, индексы, транзакции и проектирование баз данных", "🗄️", 3));
        mockTopics.add(createTopic("algorithms", "Алгоритмы",
                "Распространённые алгоритмы, анализ сложности и реализация структур данных", "⚙️", 4));
        mockTopics.add(createTopic("spring-boot", "Spring Boot",
                "Фреймворк Spring, внедрение зависимостей, веб‑разработка и микросервисы", "🍃", 5));
        mockTopics.add(createTopic("system-design", "Проектирование систем",
                "Проектирование масштабируемых систем, архитектурные шаблоны и компромиссы", "🏗️", 6));
        mockTopics.add(createTopic("devops", "DevOps и облака",
                "CI/CD, контейнеризация, облачные платформы и инфраструктура как код", "🚀", 7));
        mockTopics.add(createTopic("javascript", "JavaScript",
                "Современный JavaScript, ES6+, браузерные API и Node.js", "📜", 8));
    }

    private void createQuestions() {
        createQuestion(
                "what-is-polymorphism",
                "What is polymorphism and what are its types in OOP?",
                """
                # Polymorphism in Object-Oriented Programming

                **Polymorphism** is one of the four fundamental principles of OOP. It allows objects of different classes to be treated as objects of a common superclass, with each object responding to the same method call in its own way.

                ## Types of Polymorphism

                ### 1. Compile-time Polymorphism (Method Overloading)
                Occurs when multiple methods in the same class share the same name but have different parameters.

                ### 2. Runtime Polymorphism (Method Overriding)
                Occurs when a subclass provides a specific implementation of a method already defined in its superclass.

                ```java
                abstract class Animal {
                    public abstract void makeSound();
                }

                class Dog extends Animal {
                    @Override
                    public void makeSound() {
                        System.out.println("Woof!");
                    }
                }

                class Cat extends Animal {
                    @Override
                    public void makeSound() {
                        System.out.println("Meow!");
                    }
                }

                Animal myPet = new Dog();
                myPet.makeSound(); // Outputs: Woof!
                ```
                """,
                Level.JUNIOR,
                "oop",
                Arrays.asList("oop", "java", "design-patterns")
        );

        createQuestion(
                "inheritance-vs-composition",
                "Inheritance vs Composition: Which one should you prefer and why?",
                """
                # Inheritance vs Composition

                Modern best practices generally favor **composition over inheritance**.

                ## Inheritance
                Establishes an "is-a" relationship. A subclass inherits all public and protected members from its parent.

                **Disadvantages:**
                - Tight coupling between parent and child classes
                - Fragile base class problem
                - Deep hierarchies become hard to maintain

                ## Composition
                Establishes a "has-a" relationship. A class contains instances of other classes as members.

                ```java
                class Engine {
                    public void start() { System.out.println("Engine starting..."); }
                }

                class Car {
                    private Engine engine;

                    public Car() { this.engine = new Engine(); }

                    public void start() { engine.start(); }
                }
                ```

                **Advantages:**
                - Loose coupling
                - Greater flexibility
                - Easier to test via dependency injection
                """,
                Level.MIDDLE,
                "oop",
                Arrays.asList("oop", "design-patterns", "java")
        );

        createQuestion(
                "solid-principles-explained",
                "Can you explain the SOLID principles with examples?",
                """
                # SOLID Principles

                SOLID is an acronym for five design principles that make software designs more understandable, flexible, and maintainable.

                - **S** — Single Responsibility Principle: A class should have only one reason to change.
                - **O** — Open/Closed Principle: Classes should be open for extension but closed for modification.
                - **L** — Liskov Substitution Principle: Objects of a superclass should be replaceable with objects of a subclass.
                - **I** — Interface Segregation Principle: No client should be forced to depend on methods it does not use.
                - **D** — Dependency Inversion Principle: High-level modules should not depend on low-level modules.
                """,
                Level.SENIOR,
                "oop",
                Arrays.asList("oop", "design-patterns", "java")
        );

        createQuestion(
                "java-memory-management",
                "How does memory management work in Java? Explain heap, stack, and garbage collection.",
                """
                # Java Memory Management

                Java's memory management provides automatic garbage collection and memory allocation.

                ## JVM Memory Structure

                - **Heap**: Where all objects are stored. Divided into Young Generation (Eden + Survivor spaces) and Old Generation.
                - **Stack**: Per-thread memory storing primitive variables, object references, and method call frames.
                - **Metaspace**: Stores class structures, static variables, and the runtime constant pool.

                ## Garbage Collection

                GC automatically removes objects that are no longer reachable. Common algorithms:
                - **G1 GC** (default since Java 9): Divides heap into regions, predictable pause times.
                - **ZGC** (Java 11+): Ultra-low latency, pause times < 10ms.

                ## Best Practices

                ```java
                // Use try-with-resources to avoid leaks
                try (FileInputStream fis = new FileInputStream("file.txt")) {
                    // Automatically closed
                }

                // Use StringBuilder in loops
                StringBuilder sb = new StringBuilder();
                for (int i = 0; i < 1000; i++) {
                    sb.append(i);
                }
                ```
                """,
                Level.MIDDLE,
                "java-core",
                Arrays.asList("java", "memory-management")
        );

        createQuestion(
                "java-concurrency-basics",
                "What are the different ways to create threads in Java? Explain with examples.",
                """
                # Thread Creation in Java

                ## 1. Implementing Runnable (preferred)

                ```java
                Runnable task = () -> System.out.println("Running in: " + Thread.currentThread().getName());
                new Thread(task).start();
                ```

                ## 2. Using ExecutorService (production-ready)

                ```java
                ExecutorService executor = Executors.newFixedThreadPool(5);
                executor.submit(() -> doWork());
                executor.shutdown();
                ```

                ## 3. CompletableFuture (async workflows)

                ```java
                CompletableFuture.supplyAsync(() -> fetchData())
                    .thenApply(data -> process(data))
                    .exceptionally(ex -> handleError(ex));
                ```

                ## 4. Virtual Threads (Java 21+)

                ```java
                try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
                    executor.submit(() -> doWork());
                }
                ```
                """,
                Level.MIDDLE,
                "java-core",
                Arrays.asList("java", "concurrency")
        );

        createQuestion(
                "sql-vs-nosql",
                "SQL vs NoSQL databases: When to use which? Explain with use cases.",
                """
                # SQL vs NoSQL

                ## SQL (Relational)
                - Fixed schema, ACID transactions, strong consistency
                - Best for: financial systems, e-commerce orders, complex queries with joins

                ## NoSQL
                - Flexible schema, horizontal scaling, various models (document, key-value, graph, column)
                - Best for: real-time analytics, social feeds, caching, unstructured data

                ## Decision Guide

                **Choose SQL if:**
                - Data integrity is critical
                - Complex joins and queries are needed
                - Schema is stable

                **Choose NoSQL if:**
                - Rapidly changing schema
                - Massive horizontal scale required
                - Unstructured or semi-structured data

                ## Polyglot Persistence
                Modern apps often combine both: SQL for transactions, Redis for caching, MongoDB for catalogs.
                """,
                Level.MIDDLE,
                "databases",
                Arrays.asList("sql", "nosql")
        );

        createQuestion(
                "rest-api-design-principles",
                "What are the key principles of REST API design?",
                """
                # REST API Design Principles

                ## 1. Resource-Based URLs (nouns, not verbs)

                ```
                GET    /api/users/123    ✅
                GET    /api/getUser?id=123  ❌
                ```

                ## 2. HTTP Methods

                - **GET**: Read
                - **POST**: Create
                - **PUT**: Full update
                - **PATCH**: Partial update
                - **DELETE**: Remove

                ## 3. Appropriate Status Codes

                - `200 OK`, `201 Created`, `204 No Content`
                - `400 Bad Request`, `401 Unauthorized`, `404 Not Found`
                - `500 Internal Server Error`

                ## 4. Filtering & Pagination

                ```
                GET /api/questions?topic=java&level=middle&page=1&limit=20
                ```
                """,
                Level.MIDDLE,
                "system-design",
                Arrays.asList("microservices", "java")
        );
    }

    private TagDto createTag(String slug, String name) {
        TagDto tag = new TagDto();
        tag.setId(idSequence.getAndIncrement());
        tag.setSlug(slug);
        tag.setName(name);
        return tag;
    }

    private TopicDto createTopic(String slug, String name, String description, String icon, int sortOrder) {
        TopicDto topic = new TopicDto();
        topic.setId(idSequence.getAndIncrement());
        topic.setSlug(slug);
        topic.setName(name);
        topic.setDescription(description);
        topic.setIcon(icon);
        topic.setSortOrder(sortOrder);
        return topic;
    }

    private TagDto findTag(String slug) {
        return mockTags.stream()
                .filter(t -> t.getSlug().equals(slug))
                .findFirst()
                .orElseThrow(() -> new IllegalArgumentException("Tag not found: " + slug));
    }

    private TopicDto findTopic(String slug) {
        return mockTopics.stream()
                .filter(t -> t.getSlug().equals(slug))
                .findFirst()
                .orElseThrow(() -> new IllegalArgumentException("Topic not found: " + slug));
    }

    private void createQuestion(String slug, String title, String answer, Level level,
                                String topicSlug, List<String> tagSlugs) {
        int id = idSequence.getAndIncrement();

        TopicDto topic = findTopic(topicSlug);
        List<TagDto> tags = tagSlugs.stream().map(this::findTag).collect(Collectors.toList());

        QuestionDetailDto detail = new QuestionDetailDto();
        detail.setId(id);
        detail.setSlug(slug);
        detail.setTitle(title);
        detail.setAnswer(answer);
        detail.setLevel(level);
        detail.setTopic(topic);
        detail.setTags(tags);

        mockQuestionsById.put(id, detail);

        QuestionListItemDto listItem = convertToListItem(detail);
        questionsByTopicSlug.computeIfAbsent(topicSlug, k -> new ArrayList<>()).add(listItem);
    }

    @Override
    public List<TopicDto> getTopics() {
        log.debug("Mock: Returning {} topics", mockTopics.size());
        return new ArrayList<>(mockTopics);
    }

    @Override
    public Optional<TopicWithQuestionsDto> getTopicBySlug(String slug, Level level) {
        log.debug("Mock: Getting topic by slug: {}, level: {}", slug, level);

        return mockTopics.stream()
                .filter(t -> t.getSlug().equals(slug))
                .findFirst()
                .map(topic -> {
                    TopicWithQuestionsDto topicWithQuestions = new TopicWithQuestionsDto();
                    topicWithQuestions.setId(topic.getId());
                    topicWithQuestions.setSlug(topic.getSlug());
                    topicWithQuestions.setName(topic.getName());
                    topicWithQuestions.setDescription(topic.getDescription());
                    topicWithQuestions.setIcon(topic.getIcon());
                    topicWithQuestions.setSortOrder(topic.getSortOrder());

                    List<QuestionListItemDto> questions = questionsByTopicSlug
                            .getOrDefault(slug, Collections.emptyList())
                            .stream()
                            .filter(q -> level == null || level.equals(q.getLevel()))
                            .collect(Collectors.toList());
                    topicWithQuestions.setQuestions(questions);

                    return topicWithQuestions;
                });
    }

    @Override
    public PaginatedQuestionsDto getQuestions(String topic, String tag, Level level, int page, int limit) {
        log.debug("Mock: Getting questions - topic: {}, tag: {}, level: {}, page: {}, limit: {}",
                topic, tag, level, page, limit);

        List<QuestionListItemDto> filtered = mockQuestionsById.values().stream()
                .map(this::convertToListItem)
                .collect(Collectors.toList());

        if (topic != null && !topic.isBlank()) {
            filtered = filtered.stream()
                    .filter(q -> q.getTopic() != null && topic.equals(q.getTopic().getSlug()))
                    .collect(Collectors.toList());
        }

        if (tag != null && !tag.isBlank()) {
            filtered = filtered.stream()
                    .filter(q -> q.getTags() != null &&
                            q.getTags().stream().anyMatch(t -> tag.equals(t.getSlug())))
                    .collect(Collectors.toList());
        }

        if (level != null) {
            filtered = filtered.stream()
                    .filter(q -> level.equals(q.getLevel()))
                    .collect(Collectors.toList());
        }

        int total = filtered.size();
        int fromIndex = Math.min((page - 1) * limit, total);
        int toIndex = Math.min(fromIndex + limit, total);

        List<QuestionListItemDto> pageData = fromIndex < toIndex ?
                filtered.subList(fromIndex, toIndex) : new ArrayList<>();

        PaginatedQuestionsDto result = new PaginatedQuestionsDto();
        result.setData(pageData);

        PaginatedQuestionsDto.PaginationDto pagination = new PaginatedQuestionsDto.PaginationDto();
        pagination.setPage(page);
        pagination.setLimit(limit);
        pagination.setTotal(total);
        result.setPagination(pagination);

        return result;
    }

    @Override
    public Optional<QuestionDetailDto> getQuestionBySlug(String slug) {
        log.debug("Mock: Getting question by slug: {}", slug);
        return mockQuestionsById.values().stream()
                .filter(q -> q.getSlug().equals(slug))
                .findFirst();
    }

    @Override
    public List<TagDto> getTags() {
        log.debug("Mock: Returning {} tags", mockTags.size());
        return new ArrayList<>(mockTags);
    }

    private QuestionListItemDto convertToListItem(QuestionDetailDto detail) {
        QuestionListItemDto item = new QuestionListItemDto();
        item.setId(detail.getId());
        item.setSlug(detail.getSlug());
        item.setTitle(detail.getTitle());
        item.setLevel(detail.getLevel());
        item.setTopic(detail.getTopic());
        item.setTags(detail.getTags());
        return item;
    }
}