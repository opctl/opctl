module.exports = {
  docs: {
    Introduction: [
      "what-why"
    ],
    Setup: [
      "setup/bare-metal",
      "setup/azure-pipelines",
      "setup/docker",
      "setup/gitlab",
      "setup/kubernetes",
      "setup/travis"
    ],
    'Zero to [Op] Hero': [
      "zero-to-hero/hello-world",
      "zero-to-hero/inputs-outputs",
      {
        type: "category",
        label: "Portable Dev",
        items: [
          "run-a-go-service",
          "run-a-react-app"
        ]
      }
    ],
    CLI: [
      {
        type: "category",
        label: "Reference",
        items: [
          "cli/reference/commands"
        ]
      }
    ],
    'Opspec': [
      "opspec/introduction",
      {
        type: "category",
        label: "Reference",
        items: [
          {
            type: "category",
            label: "Structure",
            items: [
              {
                type: "category",
                label: "Op [directory]",
                items: [
                  "opspec/reference/structure/op-directory/index",
                  {
                    type: "category",
                    label: "Op [object]",
                    items: [
                      "opspec/reference/structure/op-directory/op/index",
                      {
                        type: "category",
                        label: "Call [object]",
                        items: [
                          "opspec/reference/structure/op-directory/op/call/index",
                          {
                            type: "category",
                            label: "Container Call [object]",
                            items: [
                              "opspec/reference/structure/op-directory/op/call/container/index",
                              "opspec/reference/structure/op-directory/op/call/container/image",
                            ]
                          },
                          "opspec/reference/structure/op-directory/op/call/loop-vars",
                          "opspec/reference/structure/op-directory/op/call/op",
                          "opspec/reference/structure/op-directory/op/call/parallel-loop",
                          "opspec/reference/structure/op-directory/op/call/predicate",
                          "opspec/reference/structure/op-directory/op/call/pull-creds",
                          "opspec/reference/structure/op-directory/op/call/rangeable-value",
                          "opspec/reference/structure/op-directory/op/call/serial-loop"
                        ]
                      },
                      {
                        type: "category",
                        label: "Parameter [object]",
                        items: [
                          "opspec/reference/structure/op-directory/op/parameter/index",
                          "opspec/reference/structure/op-directory/op/parameter/array",
                          "opspec/reference/structure/op-directory/op/parameter/boolean",
                          "opspec/reference/structure/op-directory/op/parameter/dir",
                          "opspec/reference/structure/op-directory/op/parameter/file",
                          "opspec/reference/structure/op-directory/op/parameter/number",
                          "opspec/reference/structure/op-directory/op/parameter/object",
                          "opspec/reference/structure/op-directory/op/parameter/socket",
                          "opspec/reference/structure/op-directory/op/parameter/string",
                        ]
                      },
                      "opspec/reference/structure/op-directory/op/initializer",
                      "opspec/reference/structure/op-directory/op/markdown",
                      "opspec/reference/structure/op-directory/op/variable-reference",
                      "opspec/reference/structure/op-directory/op/variable-name"
                    ]
                  }
                ]
              }
            ]
          },
          {
            type: "category",
            label: "Types",
            items: [
              "opspec/reference/types/array",
              "opspec/reference/types/boolean",
              "opspec/reference/types/dir",
              "opspec/reference/types/file",
              "opspec/reference/types/number",
              "opspec/reference/types/object",
              "opspec/reference/types/socket",
              "opspec/reference/types/string",
            ]
          }
        ]
      }
    ],
    SDKs: [
      {
        type: "link",
        label: "Go",
        href: "https://github.com/opctl/opctl/tree/master/sdks/go"
      },
      {
        type: "link",
        label: "Js",
        href: "https://github.com/opctl/opctl/tree/master/sdks/js"
      },
      {
        type: "link",
        label: "React",
        href: "https://github.com/opctl/opctl/tree/master/sdks/react"
      }
    ]
  }
};
