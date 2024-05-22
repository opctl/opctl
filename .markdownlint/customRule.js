/** @type import("markdownlint").Rule */
module.exports = {
  "names": [ "CustomRule/keep-a-changelog-format" ],
  "description": "CHANGELOG.md must follow the keep-a-changelog format.",
  "information": new URL("https://keepachangelog.com/en/1.1.0/"),
  "tags": [ "test", "lint", "changelog" ],
  "parser": "markdownit",
  "function": function rule(params, onError) {
    params.parsers.markdownit.tokens.filter(function filterToken(token) {
      return token.type === "heading_open";
    }).forEach(function forToken(heading, index, headings) {
      if (heading.tag === "h1" && heading.line !== "# Change Log") {
        // this explicit text is the only valid content for the first heading
        return onError({
          "lineNumber": heading.lineNumber,
          "detail": "First heading should be '# Change Log'.",
        });
      } else if (heading.tag === "h2" && !/^## \d\.\d{1,9}\.\d{1,3}-?[A-Za-z0-9]{0,} - \d{4}-\d{2}-\d{2}$/.test(heading.line)) {
        // every second heading should be a version number, then a dash and a space, and then a date
        return onError({
          "lineNumber": heading.lineNumber,
          "detail": "Second heading should be a version number and date.",
        });
      } else if (heading.tag === "h3") {
        // every third heading should contain some descriptive information about the changes in the version being described
        // and should be preceded either by a version heading or another change type heading
        if (!/^### (Added|Changed|Deprecated|Removed|Fixed)$/.test(heading.line)) {
          return onError({
            "lineNumber": heading.lineNumber,
            "detail": "Third heading should be a change type.",
          });
        }

        if (!["h3", "h2"].includes(headings[index - 1].tag)) {
          return onError({
            "lineNumber": heading.lineNumber,
            "detail": "Change type should be preceded by a version number and date.",
          });
        }
      }
    });
  }
};
