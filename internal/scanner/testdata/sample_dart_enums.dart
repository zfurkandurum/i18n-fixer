// Generic Dart fixture exercising the new hardcoded-detection patterns:
// switch arrows, return statements, throws, and additional named params.
// Strings prefixed with "FOUND_" should be flagged; "OK_" must be ignored.

enum TaskStatus {
  todo,
  doing,
  done;

  // Arrow-syntax switch — previously missed.
  String get displayName => switch (this) {
        TaskStatus.todo => 'FOUND_To Do',
        TaskStatus.doing => 'FOUND_In Progress',
        TaskStatus.done => 'FOUND_Done',
      };

  // Translation key passed to .tr() — must NOT be flagged
  // (dotted-identifier exclusion).
  String get localized => 'OK_tasks.status.$name'.tr();

  // Interpolated value — must NOT be flagged.
  String describe() => 'OK_value $name';
}

enum NotificationLevel {
  info,
  warning;

  // Classic switch/case with return — previously missed.
  String get displayName {
    switch (this) {
      case NotificationLevel.info:
        return 'FOUND_Information';
      case NotificationLevel.warning:
        return 'FOUND_Warning';
    }
  }
}

class Sample {
  // Throw with literal — should be flagged.
  void doThing() {
    throw 'FOUND_Something failed';
  }

  // Exception constructor — should be flagged.
  void other() {
    throw Exception('FOUND_Bad input');
  }

  // Named-parameter widget calls — newly covered.
  Widget render() => Column(children: [
        Text('FOUND_Hello world'),
        TextField(
          helperText: 'FOUND_Enter your name',
          errorText: 'FOUND_Invalid email',
          placeholder: 'FOUND_e.g. Jane Doe',
        ),
        ChoiceChip(label: 'FOUND_All Items'),
        ListTile(subtitle: 'FOUND_Last seen today'),
        EmptyState(actionText: 'FOUND_Add Entry'),
        InfoCard(description: 'FOUND_Brief explanation'),
        AppButton(text: 'FOUND_Submit'),
      ]);

  // Already translated via .tr() — must NOT be flagged
  // (key includes a dot).
  String localized() => 'OK_common.cancel'.tr();
}
