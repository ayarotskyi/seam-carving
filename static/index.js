window.addEventListener("load", () => {
  let img;
  const dropbox = document.getElementById("dropbox");
  document.body.addEventListener("dragenter", dragenter, false);
  document.body.addEventListener("dragover", dragover, false);
  document.body.addEventListener("drop", drop, false);
  const preview = document.getElementById("preview");

  const form = document.querySelector("form");
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
    img.onload = function () {
      document
        .getElementById("image-size")
        .replaceChildren(
          "Image size: " + this.naturalWidth + "x" + this.naturalHeight
        );
      // alert(this.naturalWidth + "x" + this.naturalHeight);
    };
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
