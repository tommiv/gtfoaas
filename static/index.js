(function() {
    const button = document.getElementById("gtfo");
    const counter = document.getElementById("counter");

    button.addEventListener("click", () => {
        const xhr = new XMLHttpRequest();
        xhr.open("POST", "/api/v1/gtfo");
        xhr.onreadystatechange = function() {
            if (xhr.readyState !== 4) return;
            const response = JSON.parse(xhr.responseText);
            if (!response.FuckedCount) return;
            counter.innerHTML = response.FuckedCount;
        }
        xhr.send();
    });
})();
