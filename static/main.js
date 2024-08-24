
const HTML_DICT= '<div class="dict" onclick="selectDict(\'%id%\')">%name%</div>';

let current_word = 0;
let words = []

function getDicts() {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', 'api/dict');
    xhr.send();
    xhr.onload = function() {
        if (xhr.status !== 200) {
            console.log(`Ошибка ${xhr.status}: ${xhr.statusText}`)
            return
        }

        let dicts = JSON.parse(xhr.response)
        let dicts_inner = '';
        for(let key in dicts) {
            let dict = dicts[key]
            let html = HTML_DICT.replace('%id%', dict.id).replace('%name%', dict.name)
            dicts_inner += html;
        }

        document.getElementById('dicts').innerHTML = dicts_inner
    };
}

function selectDict(dict_id) {
    let xhr = new XMLHttpRequest();
    xhr.open('GET', 'api/dict/' + dict_id);
    xhr.send();
    xhr.onload = function() {
        if (xhr.status !== 200) {
            console.log(`Ошибка ${xhr.status}: ${xhr.statusText}`)
            return
        }

        words = JSON.parse(xhr.response)
        nextWord();
        document.getElementById('dicts').style.display = 'none'
        document.getElementById('word').style.display = 'block'
    };
}

function nextWord() {
    show(1)
    document.getElementById('first').innerText = words[current_word]['first']
    document.getElementById('second').innerText = words[current_word]['second']

    current_word++
    if (words.length === current_word) {
        current_word = 0
    }
}

function show(toggle) {
    if (toggle) {
        document.getElementById('first').style.display = 'block'
        document.getElementById('second').style.display = 'none'
    } else {
        document.getElementById('first').style.display = 'none'
        document.getElementById('second').style.display = 'block'
    }
}

getDicts();