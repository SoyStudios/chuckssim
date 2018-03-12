const overlay = document.getElementById('overlay');
const select = document.getElementById('bot-edit-mode');
const botNameElement = document.getElementById('bot-name');
const botCodeElement = document.getElementById('bot-code');

toggleEditForm = () => {
    overlay.style.display = getComputedStyle(overlay, null).display === 'none' ? 'block' :'none';
    select.innerHTML = '';
    const defaultOption = document.createElement('option');
    const defaultText = document.createTextNode('Create new Bot');
    defaultOption.appendChild(defaultText);
    select.appendChild(defaultOption);
    for (let i of Object.keys(localStorage)) {
        const option = document.createElement('option');
        option.setAttribute('value', i);
        const t = document.createTextNode(i);
        option.appendChild(t);
        select.appendChild(option);
    }
};

fillEditorWithSelectedBot = () => {
    const selectedBot = select.options[select.selectedIndex].value;
    const code = localStorage.getItem(selectedBot);
    botNameElement.value = selectedBot;
    botCodeElement.value = code;
};

saveBot = () => {
    const name = botNameElement.value;
    const code = botCodeElement.value;
    localStorage.setItem(name, code);
};