<!DOCTYPE html><html lang="en" class="no-js">
  <head>
    <META http-equiv="Content-Type" content="text/html; charset=utf-8">
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{ .Title }}</title><script type="text/javascript" src="../res/js/jquery.js"></script>
    <script type="text/javascript" src="../res/js/toc.js"></script>
    <script type="text/javascript" src="../res/js/core.js"></script>
    <link rel="stylesheet" type="text/css" href="../res/css/core.css">
  </head>
  <body>
    <div id="master-wrapper" class="master-wrapper clear">
      <div id="sidebar" class="sidebar">
        <div class="sidebar-wrap">
          <div class="content">
            <div class="sidebar-menu">
              <div class="toc">
                <h2>Table of Contents</h2>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div id="content-wrap" class="content-wrap">
        <div class="content-block">
          <div class="content">
            <div class="section">
                {{ .Body }}
            </div>
            <div class="footer-wrapper">
              <div class="footer clear">
                <div class="copy">© Copy Right Of Your Company</div>
                <div class="menu"><a href="https://github.com/xophiix" target="_blank">xophiix</a></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </body>
</html>