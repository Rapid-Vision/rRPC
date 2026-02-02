import re
from dataclasses import dataclass
from typing import Dict, List

from rpcserver import ValidationRPCError, create_app
from rpcserver.handlers import RPCHandlers
from rpcserver.models import SliceModel, StatsModel, TextModel


@dataclass
class TextEntry:
    text: TextModel


class Service(RPCHandlers):
    def __init__(self) -> None:
        self._next_id = 1
        self._store: Dict[int, TextEntry] = {}

    def submit_text(self, text: TextModel) -> int:
        text_id = self._next_id
        self._next_id += 1
        self._store[text_id] = TextEntry(text=text)
        return text_id

    def compute_stats(self, text_id: int) -> StatsModel:
        entry = self._store.get(text_id)
        if entry is None:
            raise ValidationRPCError("unknown text id")
        return build_stats(entry.text.data)


def build_stats(data: str) -> StatsModel:
    words = re.findall(r"\\S+", data)
    word_count: Dict[str, int] = {}
    for word in words:
        word_count[word] = word_count.get(word, 0) + 1

    sentences = split_sentences(data)
    return StatsModel(
        ascii=is_ascii(data),
        word_count=word_count,
        total_words=len(words),
        sentences=sentences,
    )


def is_ascii(value: str) -> bool:
    return all(ord(ch) < 128 for ch in value)


def split_sentences(data: str) -> List[SliceModel]:
    slices: List[SliceModel] = []
    start = 0
    for idx, ch in enumerate(data):
        if ch in ".!?":
            if idx + 1 > start:
                slices.append(SliceModel(begin=start, end=idx + 1))
            start = idx + 1
    if start < len(data):
        slices.append(SliceModel(begin=start, end=len(data)))
    return slices


app = create_app(Service())

if __name__ == "__main__":
    import uvicorn

    uvicorn.run("server:app", host="127.0.0.1", port=8080, reload=False)
