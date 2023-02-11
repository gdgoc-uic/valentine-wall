<html>
    <body>
        <script type="text/javascript">
            const params = Object.fromEntries((new URL(window.location)).searchParams.entries());
            window.opener.postMessage(params, {{ .FrontendURL }});
        </script>
    </body>
</html>