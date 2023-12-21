import subprocess
import tempfile
from typing import Final, Optional
from pathlib import Path


class BlogManager:
    def __init__(self, repo_url: str, repo_path: Optional[str] = None) -> None:
        self.repo_url: Final[str] = repo_url
        self.repo: str = repo_path or self.get_repo()


    def get_repo(self) -> str:
        repo = tempfile.mkdtemp()
        cmd = ["git", "clone", self.repo_url, repo]
        subprocess.run(cmd, timeout=10)
        return repo

    @property
    def articles(self):
        files: list[Path] = []
        self.__add_folder(Path(self.repo), files)
        return files

    def __add_folder(self, folder: Path, files: list[Path]) -> list[Path]:
        for f in folder.iterdir():
            if f.is_file() and f.suffix == ".md":
                files.append(f)

            elif f.is_dir() and not f.name.startswith('.'):
                self.__add_folder(f, files)

        return files
