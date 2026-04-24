# CodeMind: Local-First AI Code

**CodeMind** es un sistema de RAG (Retrieval-Augmented Generation) de alto rendimiento construido en **Go** y **Wails**. Permite indexar repositorios de código locales y realizar consultas en lenguaje natural sobre la arquitectura, lógica y flujo del código, garantizando privacidad total al procesar todo mediante **Ollama** de forma local.

## 🚀 Arquitectura del Sistema

El proyecto se divide en tres capas principales que aprovechan el hardware local (GPU/VRAM) y la concurrencia de Go.

### 1. Ingestion Pipeline (The "Indexer")
Escrito en Go para máxima velocidad.
- **File Crawler:** Escanea el directorio ignorando archivos en `.gitignore`.
- **AST-Based Chunking:** Utiliza `tree-sitter-go` para fragmentar el código por funciones, métodos y estructuras en lugar de líneas de texto.
- **Embedding Generator:** Interface con la API local de Ollama para convertir fragmentos de código en vectores usando el modelo `mxbai-embed-large`.

### 2. Vector Storage
- **Engine:** Qdrant (corriendo en Docker) o una implementación embebida en Go como `BadgerDB` con búsqueda de similitud de coseno.
- **Schema:** Almacena el vector + metadatos (file_path, line_start, language, symbols).

### 3. Retrieval & Augmentation (The "Brain")
- **Hybrid Search:** Combina búsqueda vectorial con búsqueda de palabras clave para precisión en nombres de variables únicos.
- **Context Window Management:** Aprovecha los 131k de contexto de modelos como Gemma 4 para enviar múltiples archivos relacionados en una sola consulta.

---

## 🛠 Tech Stack

- **Lenguaje:** Go 1.26 (Backend)
- **Frontend:** React + Tailwind CSS (vía Wails)
- **IA Local:** Ollama (Models: Gemma 4 para razonamiento, mxbai-embed-large para vectores de preferencia esos pero pueden ser configurables con lo que tenga ollama)
- **Base de Datos:** Qdrant / BadgerDB
- **Parsing:** Tree-sitter (para análisis sintáctico de código)
- **Estructura de Carpetas:**
    /cmd/run.go
    /internal
        /engine            # Lógica central del RAG (Orquestador)
        /ingestion         # Crawler, Chunker y Embeddings
        /storage           # Implementación de Vector DB (Qdrant/Local)
        /api               # Handlers para Wails / REST
    /pkg                   # Utilidades agnósticas (Parser, Logger)
    main.go (Usar cobra)

---

## 📋 Hoja de Ruta de Desarrollo

### Fase 1: Core de Ingestión (CLI en Go)
- [ ] Configurar el cliente de Ollama en Go para generar embeddings.
- [ ] Implementar el "Chunker" inteligente que reconozca bloques de código completos.
- [ ] Script de migración masiva a la Base de Datos Vectorial.

### Fase 2: Integración de IA & RAG
- [ ] Implementar la lógica de búsqueda de similitud (Similarity Search).
- [ ] Diseñar el "System Prompt" para que el LLM actúe como un Senior Architect.
- [ ] Crear el flujo de orquestación: Pregunta -> Vector -> Búsqueda -> Prompt -> Respuesta.

### Fase 3: Interfaz de Usuario (Wails)
- [ ] Dashboard de indexación (visualizar progreso y archivos procesados).
- [ ] Interfaz de chat tipo VS Code.
- [ ] Visualizador de fuentes (mostrar qué fragmentos de código usó la IA para dar la respuesta).

---

## 💡 Key Performance Indicators (Para el Portafolio)

Para demostrar tu seniority, el proyecto incluirá métricas de:
- **Indexing Speed:** Tiempo de procesamiento por cada 1,000 líneas de código (configurable).
- **Retrieval Latency:** Tiempo de respuesta desde la pregunta hasta la obtención del contexto.
- **VRAM Optimization:** Monitorización de carga de GPU durante la inferencia local (mostrar en la UI).

---

## 🛡 Privacidad & Seguridad
- **Zero Cloud:** Ningún dato sale de la red local.
- **Local LLM:** Soporte para modelos descargados vía Ollama.
- **Analytic Free:** Sin telemetría externa.