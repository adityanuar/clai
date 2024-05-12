from fastapi import FastAPI
from pydantic import BaseModel
from llmware.prompts import Prompt

model_name = "llmware/dragon-yi-6b-gguf"
prompter = Prompt().load_model(model_name)

app = FastAPI()

class Prompt(BaseModel):
    content: str

@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.post("/suggest/")
def suggest(prompt: Prompt):
    response = prompter.prompt_main(str(prompt.content) + "? answer with actionable commands only without description", context="", prompt_name="default_with_context", temperature=0.9)
    return {"result": f"{response['llm_response']}"}