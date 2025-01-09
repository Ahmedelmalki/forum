    function isEmptyOrWhitespace(input) {
      return !input.trim();
    }
    document.getElementById('newPostForm').addEventListener('submit', function (event) {
      const title = document.getElementById('title').value;
      const content = document.getElementById('content').value;
      const categories = Array.from(document.querySelectorAll('input[name="categories[]"]:checked')).map(input => input.value);
      if (isEmptyOrWhitespace(title)) {
        event.preventDefault(); 
        alert("Post title shouldn't be empty or only whitespace.");
        return;
      }
      if (isEmptyOrWhitespace(content)) {
        event.preventDefault(); 
        alert("Post content shouldn't be empty or only whitespace.");
        return;
      }
      if (categories.length === 0) {
        event.preventDefault(); 
        alert("Please select at least one category.");
        return;
      }
      const createPostButton = document.getElementById('createPostButton');
      createPostButton.disabled = true;
      createPostButton.textContent = 'Submitting...';

      setTimeout(() => {
        createPostButton.disabled = false;
        createPostButton.textContent = 'Create Post';
      }, 5000);
    });
