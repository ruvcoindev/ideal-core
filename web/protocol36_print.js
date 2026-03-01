// protocol36_print.js — клиентский скрипт для печати PDF

// Генерация PDF через браузерный print
function generateProtocolPDF(personName, lastContact, phaseConfig) {
	const printContent = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>Protocol36: ${personName}</title>
	<style>
		body { font-family: monospace; padding: 20px; max-width: 800px; margin: 0 auto; }
		.header { text-align: center; border-bottom: 2px solid #333; padding-bottom: 10px; margin-bottom: 20px; }
		.phase { background: #f5f5f5; padding: 15px; margin: 10px 0; border-left: 4px solid #00d9ff; }
		.task { margin: 5px 0 5px 20px; }
		.warning { background: #fff3cd; padding: 10px; border: 1px solid #ffc107; margin: 10px 0; }
		.footer { margin-top: 30px; font-size: 0.9em; color: #666; }
	</style>
</head>
<body>
	<div class="header">
		<h1>🗝️ Protocol36: ${personName}</h1>
		<p><strong>Старт:</strong> ${lastContact} | <strong>Завершение:</strong> ${calculateEndDate(lastContact)}</p>
	</div>

	<div class="warning">
		⚠️ Это не замена терапии. При кризисе: 112 или 8-800-2000-122
	</div>

	<h2>📅 Фазы протокола</h2>
	${generatePhaseBlocks(phaseConfig)}

	<h2>🧘 Ежедневные практики</h2>
	${generateTaskList(phaseConfig.DailyTasks)}

	<div class="footer">
		<p>Создано в ideal-core | ${new Date().toLocaleDateString('ru-RU')}</p>
		<p>Ваша ценность не зависит от выбора другого человека.</p>
	</div>

	<script>
		window.onload = () => {
			setTimeout(() => window.print(), 500);
		};
		function calculateEndDate(start) {
			const d = new Date(start);
			d.setDate(d.getDate() + 36);
			return d.toLocaleDateString('ru-RU');
		}
		function generatePhaseBlocks(config) {
			return \`
				<div class="phase">
					<strong>🌑 Detox (Дни 1-7)</strong><br>
					Фокус: Тишина, наблюдение, заземление
				</div>
				<div class="phase">
					<strong>🌓 Rewire (Дни 8-21)</strong><br>
					Фокус: Новые практики, мягкие эксперименты
				</div>
				<div class="phase">
					<strong>🌕 Integration (Дни 22-36)</strong><br>
					Фокус: Закрепление, тестирование в реальности
				</div>
			\`;
		}
		function generateTaskList(tasks) {
			return tasks.map(t => \`<div class="task">• <strong>\${t.Title}</strong>: \${t.Description}</div>\`).join('');
		}
	</script>
</body>
</html>`;

	const printWindow = window.open('', '_blank');
	printWindow.document.write(printContent);
	printWindow.document.close();
}

// Экспорт в JSON для бэкапа
function exportProtocolJSON(personName, lastContact, phaseConfig) {
	const data = {
		person: personName,
		startDate: lastContact,
		endDate: new Date(lastContact).setDate(new Date(lastContact).getDate() + 36),
		phaseConfig: phaseConfig
	};
	const blob = new Blob([JSON.stringify(data, null, 2)], {type: 'application/json'});
	const url = URL.createObjectURL(blob);
	const a = document.createElement('a');
	a.href = url;
	a.download = `protocol36_${personName}.json`;
	a.click();
}
