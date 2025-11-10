package scripts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmbeddedScripts(t *testing.T) {
	t.Run("localstack init script is embedded", func(t *testing.T) {
		assert.NotEmpty(t, LocalstackInitScript)
		assert.Contains(t, LocalstackInitScript, "#!/")
		assert.Contains(t, strings.ToLower(LocalstackInitScript), "localstack")
	})

	t.Run("kafka topics init script is embedded", func(t *testing.T) {
		assert.NotEmpty(t, KafkaTopicsInitScript)
		assert.Contains(t, KafkaTopicsInitScript, "#!/")
		assert.Contains(t, strings.ToLower(KafkaTopicsInitScript), "kafka")
	})

	t.Run("scripts have valid shell headers", func(t *testing.T) {
		scripts := map[string]string{
			"LocalstackInitScript":  LocalstackInitScript,
			"KafkaTopicsInitScript": KafkaTopicsInitScript,
		}

		for name, script := range scripts {
			assert.True(t, strings.HasPrefix(script, "#!/"),
				"Script %s should start with shebang", name)
		}
	})

	t.Run("scripts are non-empty and contain content", func(t *testing.T) {
		scripts := map[string]string{
			"LocalstackInitScript":  LocalstackInitScript,
			"KafkaTopicsInitScript": KafkaTopicsInitScript,
		}

		for name, script := range scripts {
			lines := strings.Split(strings.TrimSpace(script), "\n")
			assert.Greater(t, len(lines), 1,
				"Script %s should have multiple lines", name)

			// Should have more than just the shebang
			nonEmptyLines := 0
			for _, line := range lines {
				if strings.TrimSpace(line) != "" && !strings.HasPrefix(strings.TrimSpace(line), "#") {
					nonEmptyLines++
				}
			}
			assert.Greater(t, nonEmptyLines, 0,
				"Script %s should have executable content", name)
		}
	})
}
