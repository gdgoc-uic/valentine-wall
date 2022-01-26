<html>
  <head>
    <title>image</title>
    <link rel="preconnect" href="https://fonts.googleapis.com">
    <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
    <link href="https://fonts.googleapis.com/css2?family=Lato&family=Nanum+Pen+Script&display=swap" rel="stylesheet">
    <style>
      * {
        box-sizing: border-box;
      }
      body {
        margin: 0;
        font-size: 16px;
      }
      p {
        margin: 0;
      }
      .image-wrapper {
        width: 1200px;
        height: 675px;
        background-color: rgb(251, 207, 232);
        padding: 3rem;
      }
      .image-wrapper .content-wrapper {
        background: linear-gradient(to bottom,rgb(254, 243, 199) 2.95rem,#00b0d7 1px); background-size: 100% 3.1rem;
        height: 100%;
        border-radius: 2rem;
        padding-left: 3rem;
        padding-right: 3rem;
        padding-top: 3rem;
        padding-bottom: 1.5rem;
        display: flex;
        flex-direction: column;
        text-overflow: ellipsis;
      }

      .image-wrapper .content-wrapper .content {
        font-family: 'Nanum Pen Script', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        line-height: 0.9;
        text-align: left;
        text-overflow: ellipsis;
        overflow: hidden;
        width: 100%;
        font-size: 3.5rem;
        height: 100%;
      }
      .image-wrapper .content-wrapper .timestamp {
        font-family: 'Lato', 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        color: rgb(10, 10, 10);
        position: relative;
        top: 1rem;
      }
    </style>
  </head>
  <body>
    <div id="image-preview" class="image-wrapper">
      <div class="content-wrapper">
        <p class="content">{{ .Content }}</p>
        <p class="timestamp">Posted on {{ .CreatedAt }}</p>
      </div>
    </div>
  </body>
</html>