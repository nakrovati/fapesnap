import { ThemeSelector } from "~/features/theme";

function SettingsPage() {
	return (
		<>
			<section>
				<h1 class="text-3xl">Settings</h1>

				<div class="mt-8">
					<p>Appearance</p>
					<ThemeSelector />
				</div>
			</section>
		</>
	);
}

export default SettingsPage;
