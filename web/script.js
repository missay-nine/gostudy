const boardSize = 15;
const emptySymbol = ".";
const humanSymbol = "X";
const aiSymbol = "O";

let board = [];
let gameOver = false;
let aiThinking = false;

const boardEl = document.getElementById("board");
const statusEl = document.getElementById("status");
const resetBtn = document.getElementById("reset");

function init() {
    document.documentElement.style.setProperty("--board-size", boardSize);
    resetBtn.addEventListener("click", resetGame);
    createBoard();
    resetGame();
}

function createBoard() {
    boardEl.style.gridTemplateColumns = `repeat(${boardSize}, var(--cell-size))`;
    boardEl.style.gridTemplateRows = `repeat(${boardSize}, var(--cell-size))`;
    boardEl.innerHTML = "";

    for (let r = 0; r < boardSize; r++) {
        for (let c = 0; c < boardSize; c++) {
            const cell = document.createElement("button");
            cell.className = "cell";
            cell.dataset.row = r;
            cell.dataset.col = c;
            cell.addEventListener("click", handleCellClick);
            cell.setAttribute("role", "gridcell");
            cell.setAttribute("aria-label", `行 ${r + 1} 列 ${c + 1}`);
            boardEl.appendChild(cell);
        }
    }
}

function resetGame() {
    board = Array.from({ length: boardSize }, () => Array(boardSize).fill(emptySymbol));
    gameOver = false;
    aiThinking = false;
    updateBoardUI();
    setStatus("轮到你落子");
}

function handleCellClick(event) {
    if (gameOver || aiThinking) return;
    const row = Number(event.currentTarget.dataset.row);
    const col = Number(event.currentTarget.dataset.col);

    if (board[row][col] !== emptySymbol) {
        setStatus("该位置已有棋子，请选择其他格子。");
        return;
    }

    placeStone(row, col, humanSymbol);
    if (checkEndState(humanSymbol)) return;

    aiThinking = true;
    setStatus("AI 思考中...");
    setTimeout(() => {
        aiMove();
        aiThinking = false;
    }, 120);
}

function aiMove() {
    const move = chooseAIMove();
    if (!move) return;
    placeStone(move.row, move.col, aiSymbol);
    if (checkEndState(aiSymbol)) return;
    setStatus("轮到你落子");
}

function placeStone(row, col, symbol) {
    board[row][col] = symbol;
    const index = row * boardSize + col;
    const cellEl = boardEl.children[index];
    cellEl.textContent = symbol;
    cellEl.classList.toggle("human", symbol === humanSymbol);
    cellEl.classList.toggle("ai", symbol === aiSymbol);
}

function checkEndState(symbol) {
    if (checkWin(symbol)) {
        gameOver = true;
        setStatus(symbol === humanSymbol ? "你赢了！" : "AI 获胜！");
        return true;
    }
    if (isDraw()) {
        gameOver = true;
        setStatus("平局！");
        return true;
    }
    return false;
}

function updateBoardUI() {
    Array.from(boardEl.children).forEach((cell) => {
        cell.textContent = "";
        cell.classList.remove("human", "ai");
    });
}

function setStatus(text) {
    statusEl.textContent = text;
}

// AI & Rules
function chooseAIMove() {
    const empties = availableCells();
    if (empties.length === 0) return null;

    const winningMove = findWinningMove(aiSymbol);
    if (winningMove) return winningMove;

    const blockingMove = findWinningMove(humanSymbol);
    if (blockingMove) return blockingMove;

    const center = Math.floor(boardSize / 2);
    if (board[center][center] === emptySymbol) {
        return { row: center, col: center };
    }

    const neighbors = prioritizedNeighbors();
    if (neighbors.length > 0) {
        return neighbors[Math.floor(Math.random() * neighbors.length)];
    }

    return empties[Math.floor(Math.random() * empties.length)];
}

function findWinningMove(symbol) {
    for (const cell of availableCells()) {
        board[cell.row][cell.col] = symbol;
        const win = checkWin(symbol);
        board[cell.row][cell.col] = emptySymbol;
        if (win) return cell;
    }
    return null;
}

function availableCells() {
    const cells = [];
    for (let r = 0; r < boardSize; r++) {
        for (let c = 0; c < boardSize; c++) {
            if (board[r][c] === emptySymbol) {
                cells.push({ row: r, col: c });
            }
        }
    }
    return cells;
}

function prioritizedNeighbors() {
    const cells = [];
    for (let r = 0; r < boardSize; r++) {
        for (let c = 0; c < boardSize; c++) {
            if (board[r][c] !== emptySymbol) continue;
            if (hasNeighbor(r, c)) {
                cells.push({ row: r, col: c });
            }
        }
    }
    return cells;
}

function hasNeighbor(row, col) {
    for (let dr = -1; dr <= 1; dr++) {
        for (let dc = -1; dc <= 1; dc++) {
            if (dr === 0 && dc === 0) continue;
            const nr = row + dr;
            const nc = col + dc;
            if (nr >= 0 && nr < boardSize && nc >= 0 && nc < boardSize) {
                if (board[nr][nc] !== emptySymbol) return true;
            }
        }
    }
    return false;
}

function isDraw() {
    return board.every((row) => row.every((cell) => cell !== emptySymbol));
}

function checkWin(symbol) {
    const directions = [
        { dr: 0, dc: 1 },
        { dr: 1, dc: 0 },
        { dr: 1, dc: 1 },
        { dr: 1, dc: -1 },
    ];

    for (let r = 0; r < boardSize; r++) {
        for (let c = 0; c < boardSize; c++) {
            if (board[r][c] !== symbol) continue;
            for (const dir of directions) {
                if (countConsecutive(r, c, dir.dr, dir.dc, symbol) >= 5) {
                    return true;
                }
            }
        }
    }
    return false;
}

function countConsecutive(row, col, dr, dc, symbol) {
    let count = 0;
    let r = row;
    let c = col;

    while (r >= 0 && r < boardSize && c >= 0 && c < boardSize && board[r][c] === symbol) {
        count++;
        r += dr;
        c += dc;
    }
    return count;
}

document.addEventListener("DOMContentLoaded", init);
