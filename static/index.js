window.addEventListener("load", () => {
  let img;
  let dropbox = document.getElementById("dropbox");
  dropbox.addEventListener("dragenter", dragenter, false);
  dropbox.addEventListener("dragover", dragover, false);
  dropbox.addEventListener("drop", drop, false);
  let preview = document.getElementById("preview");

  let form = document.querySelector("form");
  form.addEventListener("formdata", (e) => {
    const formData = e.formData;

    formData.append("image", img.file);
  });

  form.addEventListener("submit", (event) => {
    event.preventDefault();
    if (!img) {
      alert("Please select an image");
      return;
    }
    fetch(event.target.action, {
      method: "POST",
      body: new FormData(event.target), // event.target is the form
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.blob();
      })
      .then((blob) => {
        preview.src = URL.createObjectURL(blob);
      });
  });

  function dragenter(e) {
    e.stopPropagation();
    e.preventDefault();
  }

  function dragover(e) {
    e.stopPropagation();
    e.preventDefault();
  }

  function drop(e) {
    e.stopPropagation();
    e.preventDefault();

    let file = e.dataTransfer.files[0];

    if (!file?.type.match(/image.*/)) {
      return;
    }

    img = document.createElement("img");
    if (FileReader) {
      var fr = new FileReader();
      fr.onload = function () {
        img.src = fr.result;
      };
      fr.readAsDataURL(file);
    }
    img.style.width = "100%";
    img.file = file;
    dropbox.replaceChildren(img);
  }
});
