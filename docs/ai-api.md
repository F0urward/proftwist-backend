# AI Service API

## Configuration

AI credentials are stored in backend config (`config/config.yml`):

```yaml
ai:
  provider: "ollama"           # default provider
  ollama:
    baseUrl: "http://host.docker.internal:11434/v1"
    apiKey: ""
    model: "phi3.5:3.8b"
  openai:
    baseUrl: "https://openrouter.ai/api/v1"
    apiKey: ""
    model: "openrouter/free"
```

Credentials can be overridden via environment variables:
- `AI_OLLAMA_BASE_URL`, `AI_OLLAMA_API_KEY`, `AI_OLLAMA_MODEL`
- `AI_OPENAI_BASE_URL`, `AI_OPENAI_API_KEY`, `AI_OPENAI_MODEL`

## Endpoints

### 1. Generate Roadmap
**POST** `/api/v1/ai/roadmap`

**Request Body:**
```json
{
  "prompt": "required - FULL prompt including all instructions for LLM",
  "roadmap_id": "optional",
  "provider": "optional",   // "ollama", "openai", "mock"
  "model": "optional"      // override default model for provider
}
```

**Response:**
Raw JSON response from LLM - frontend controls the output format entirely.

---

### 2. Generate Roadmap Node Description
**POST** `/api/v1/ai/roadmap-node-description`

**Request Body:**
```json
{
  "node_id": "required",
  "node_label": "required",
  "roadmap_id": "optional",
  "node_type": "optional",
  "current_description": "optional",
  "provider": "optional",
  "model": "optional"
}
```

**Response:**
```json
{
  "description": "Generated description"
}
```

---

## How It Works

The `/api/v1/ai/roadmap` endpoint is now a **passthrough** - it sends the exact prompt from frontend to the LLM and returns the raw LLM response without any backend transformation.

Frontend controls:
- Full prompt content (system instructions, format requirements, etc.)
- Output format specification
- All LLM parameters via the prompt

Backend provides:
- Provider selection (ollama/openai/mock)
- Model selection
- Authentication/connection to LLM

## Example Request

```json
{
  "prompt": "You are an expert in creating developer roadmaps. Generate a logical roadmap based on the user's request.\n\nSTRICT RULES:\n\n1. Output ONLY raw JSON...\n2. JSON keys MUST be in English...\n3. Values for label and description MUST be in Russian...\n\nUSER REQUEST:\nReact",
  "provider": "openai",
  "model": "openrouter/auto"
}
```

The backend will send this exact prompt to the LLM and return whatever the LLM responds with - no transformation or modification.