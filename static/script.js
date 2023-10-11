document.getElementById("start-button").addEventListener("click", function () {
  var resultDiv = document.getElementById("result");
  resultDiv.innerHTML = "Testing...";

  fetch("/speedtest")
    .then((response) => response.text())
    .then((data) => {
      resultDiv.innerHTML = data;
    })
    .catch((error) => {
      resultDiv.innerHTML = "Error: " + error;
    });
});
