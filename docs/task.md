## TaskScheduler

The `TaskScheduler` is a fundamental component of the gAgents framework, allowing you to schedule and execute tasks at specific times or intervals.

### How it Works

The `TaskScheduler` manages a collection of tasks, each identified by a unique ID. Tasks are added to the scheduler along with their execution schedule.

### Adding a Task

To add a task, you create a struct that implements the `Task` interface. This interface defines the necessary methods for a task, including `ID()`, `Type()`, `ScheduledAt()`, `Parameters()`, `Execute()`, `Interval()`, and `StopCondition()`.

### Executing Tasks

The scheduler continuously checks the current time against the scheduled execution time of each task. If the scheduled time has passed, the task's `Execute()` method is called.

### Recurring Tasks

Tasks can be set to recur at specific intervals. If a task has a non-zero interval defined in its `Interval()` method, the scheduler will reschedule the task for the next occurrence after execution.

### Stopping Tasks

Tasks may have a stop condition defined in their `StopCondition()` method. If this condition is met after execution, the task is removed from the scheduler.

### Example Usage

```go
// Create a new context for the scheduler
ctx, cancel := context.WithCancel(context.Background())
scheduler := gAgents.NewTaskScheduler(ctx)

// Define a task and add it to the scheduler
myTask := MyCustomTask{
    // Task properties...
}

scheduler.AddTask(myTask)

// Start executing tasks
go scheduler.ExecuteTasks()
```