window.addEventListener("load", () => {
  const text = document.getElementsByTagName("h1").item(0);

  if (!text) {
    return;
  }

  text.animate(
    [
      { transform: "scale(1)" },
      { transform: "scale(0.5)" },
      { transform: "scale(1)" },
    ],
    {
      duration: 3000,
      iterations: Infinity,
    }
  );

  fetch("/hello")
    .then((response) => response.json())
    .then(console.log);
});
