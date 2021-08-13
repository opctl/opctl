package opgraph

import (
	"errors"
	"testing"
	"time"

	. "github.com/onsi/gomega"
	"github.com/opctl/opctl/sdks/go/model"
)

type noopOpFormatter struct{}

func (noopOpFormatter) FormatOpRef(opRef string) string {
	return opRef
}

func TestCallGraph(t *testing.T) {
	g := NewGomegaWithT(t)

	/* arrange */
	timestamp, err := time.Parse("Jan 2, 2006 at 3:04pm (MST)", "Feb 4, 2014 at 6:05pm (PST)")
	if err != nil {
		t.Fatal(err)
	}
	objectUnderTest := CallGraph{}
	parentID := "parentID"
	child1ID := "child1ID"
	child2ID := "child2ID"
	child3ID := "child3ID"
	containerRef := "containerRef"
	child1child1ID := "child1child1Id"
	child2If := false
	child3If := true

	g.Expect(objectUnderTest.String(
		StaticLoadingSpinner{},
		timestamp.Add(time.Second*60),
		true,
	)).To(Equal("Empty call graph"))

	/* act */
	// this is the root event
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID: parentID,
				Op: &model.OpCall{
					BaseCall: model.BaseCall{
						OpPath: "oppath",
					},
					OpID: "firstopid",
				},
			},
		},
		Timestamp: timestamp,
	})
	// first child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child1ID,
				ParentID: &parentID,
				Op: &model.OpCall{
					BaseCall: model.BaseCall{
						OpPath: "oppath2",
					},
					OpID: "secondopid",
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 1),
	})
	// first child -> first child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child1child1ID,
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "id1234567890",
					Cmd:         []string{"./cmd", "arg1", "arg2", "arg3"},
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 2),
	})
	// first child -> second child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child1Child2Id",
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 3),
	})
	// first child -> second child succeeds
	objectUnderTest.HandleEvent(&model.Event{
		CallEnded: &model.CallEnded{
			Call: model.Call{
				ID:       "child1Child2Id",
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
			Outcome: model.OpOutcomeSucceeded,
		},
		Timestamp: timestamp.Add(time.Second * 4),
	})
	// second child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child2ID",
				ParentID: &parentID,
				If:       &child2If,
			},
		},
		Timestamp: timestamp.Add(time.Second * 5),
	})
	// third child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child3ID",
				ParentID: &parentID,
				If:       &child3If,
				Serial:   []*model.CallSpec{},
			},
		},
		Timestamp: timestamp.Add(time.Second * 6),
	})
	// first child -> third child
	child1Child3ID := "child1Child3Id"
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child1Child3ID,
				ParentID: &child1ID,
			},
		},
		Timestamp: timestamp.Add(time.Second * 3),
	})
	// first child -> third child -> first child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child1Child3Child1Id",
				ParentID: &child1Child3ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 3),
	})
	// first child -> third child failed
	objectUnderTest.HandleEvent(&model.Event{
		CallEnded: &model.CallEnded{
			Call: model.Call{
				ID:       "child1Child3Id",
				ParentID: &child1ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
			Outcome: model.OpOutcomeFailed,
		},
		Timestamp: timestamp.Add(time.Second * 4),
	})
	// second child -> first child
	containerName := "named-container"
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child2Child1Id",
				ParentID: &child2ID,
				Container: &model.ContainerCall{
					Name:        &containerName,
					ContainerID: "id0987654321",
					Cmd:         []string{"./cmd"},
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 3),
	})
	// second child -> first child killed
	objectUnderTest.HandleEvent(&model.Event{
		CallEnded: &model.CallEnded{
			Call: model.Call{
				ID:       "child2Child1Id",
				ParentID: &child2ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
			Outcome: model.OpOutcomeKilled,
		},
		Timestamp: timestamp.Add(time.Second * 4),
	})
	// third child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child3ID,
				ParentID: &parentID,
			},
		},
		Timestamp: timestamp.Add(time.Second * 1),
	})
	// third child -> first child
	thirdChildName := "Third Child"
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				Name:     &thirdChildName,
				ID:       "child3IDChild1ID",
				ParentID: &child3ID,
				Parallel: []*model.CallSpec{},
			},
		},
		Timestamp: timestamp.Add(time.Second * 1),
	})
	// third child succeeds, so its children should be collapsed
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child3ID,
				ParentID: &parentID,
			},
		},
		Timestamp: timestamp.Add(time.Second * 9),
	})
	// fourth child
	child4ID := "child4ID"
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child4ID,
				ParentID: &parentID,
				Op: &model.OpCall{
					BaseCall: model.BaseCall{
						OpPath: "oppath2",
					},
					OpID: "secondopid",
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 1),
	})
	// fourth child -> first child
	child4child1ID := "child4child1ID"
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       child4child1ID,
				ParentID: &child4ID,
				Container: &model.ContainerCall{
					ContainerID: "id1234567890",
					Cmd:         []string{"./cmd", "arg1", "arg2", "arg3"},
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 2),
	})
	// fourth child -> second child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child4Child2Id",
				ParentID: &child4ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
		},
		Timestamp: timestamp.Add(time.Second * 3),
	})
	// fourth child -> second child succeeds
	objectUnderTest.HandleEvent(&model.Event{
		CallEnded: &model.CallEnded{
			Call: model.Call{
				ID:       "child4Child2Id",
				ParentID: &child4ID,
				Container: &model.ContainerCall{
					ContainerID: "id0987654321",
					Image: &model.ContainerCallImage{
						Ref: &containerRef,
					},
				},
			},
			Outcome: model.OpOutcomeSucceeded,
		},
		Timestamp: timestamp.Add(time.Second * 4),
	})
	// fourth child succeeds (should be collapsed)
	objectUnderTest.HandleEvent(&model.Event{
		CallEnded: &model.CallEnded{
			Call: model.Call{
				ID:       child4ID,
				ParentID: &parentID,
			},
			Outcome: model.OpOutcomeSucceeded,
		},
		Timestamp: timestamp.Add(time.Second * 30),
	})
	// fifth child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:           "child5ID",
				ParentID:     &parentID,
				ParallelLoop: &model.ParallelLoopCall{},
			},
		},
		Timestamp: timestamp.Add(time.Second * 9),
	})
	// sixth child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:       "child6ID",
				ParentID: &parentID,
				Serial:   []*model.CallSpec{},
			},
		},
		Timestamp: timestamp.Add(time.Second * 9),
	})
	// seventh child
	objectUnderTest.HandleEvent(&model.Event{
		CallStarted: &model.CallStarted{
			Call: model.Call{
				ID:         "child7ID",
				ParentID:   &parentID,
				SerialLoop: &model.SerialLoopCall{},
			},
		},
		Timestamp: timestamp.Add(time.Second * 9),
	})

	objectUnderTest.errors = append(objectUnderTest.errors, errors.New("this should show up as a warning"))

	/* assert */
	// the newline is here just for better test code readability
	collapsedStr := "\n" + objectUnderTest.String(
		StaticLoadingSpinner{},
		timestamp.Add(time.Second*60),
		true,
	)
	expectedCollapsedStr := `
◎ oppath
├─◎ oppath2
│ ├─◉ ⋰ id123456 containerRef 58s ./cmd arg1 arg2 arg3
│ ├─◉ ☑ id098765 containerRef 1s
│ └─◎ ⚠ 1s
│   └─◉ ⋰ id098765 containerRef 57s
├─◎ if skipped
│ └─◉️ ☒ id098765 named-container 1s
├─◎ if
│ │ serial
│ └─◉ ⋰ Third Child parallel 59s
├─◉ ⋰ 59s
├─◉ ⋰ 51s
├─◎ ☑ oppath2 29s (2 children)
├─◉ ⋰ parallel loop 51s
├─◉ ⋰ serial 51s
└─◉ ⋰ serial loop 51s
⚠️  this should show up as a warning`
	g.Expect(collapsedStr).To(Equal(expectedCollapsedStr))

	// the newline is here just for better test code readability
	expandedStr := "\n" + objectUnderTest.String(
		StaticLoadingSpinner{},
		timestamp.Add(time.Second*60),
		false,
	)
	expandedCollapsedStr := `
◎ oppath
├─◎ oppath2
│ ├─◉ ⋰ id123456 containerRef 58s ./cmd arg1 arg2 arg3
│ ├─◉ ☑ id098765 containerRef 1s
│ └─◎ ⚠ 1s
│   └─◉ ⋰ id098765 containerRef 57s
├─◎ if skipped
│ └─◉️ ☒ id098765 named-container 1s
├─◎ if
│ │ serial
│ └─◉ ⋰ Third Child parallel 59s
├─◉ ⋰ 59s
├─◉ ⋰ 51s
├─◎ ☑ oppath2 29s
│ ├─◉ ⋰ id123456 containerRef 58s ./cmd arg1 arg2 arg3
│ └─◉ ☑ id098765 containerRef 1s
├─◉ ⋰ parallel loop 51s
├─◉ ⋰ serial 51s
└─◉ ⋰ serial loop 51s
⚠️  this should show up as a warning`
	g.Expect(expandedStr).To(Equal(expandedCollapsedStr))
}
