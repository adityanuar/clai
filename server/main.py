from fastapi import FastAPI
from pydantic import BaseModel
from llmware.prompts import Prompt
import re

# Load the model, can be:
# 1. llmware/dragon-yi-6b-gguf
# 2. llmware/bling-phi-3-gguf
# etc
model_name = "llmware/dragon-yi-6b-gguf"
prompter = Prompt().load_model(model_name)

app = FastAPI()

class Prompt(BaseModel):
    content: str
    os: str

@app.get("/")
def read_root():
    return {"Hello": "World"}


@app.post("/suggest/")
def suggest(prompt: Prompt):
    print(prompt.os)
    response = prompter.prompt_main('shell command to ' + str(prompt.content) + "? answer with first one actionable commands only without description", context=f"operating system is {prompt.os}", prompt_name="default_with_context", temperature=0.9)
    filtered_response = re.sub(r"^\d+\.\s*", "", response['llm_response'])
    return {"result": f"{filtered_response}"}