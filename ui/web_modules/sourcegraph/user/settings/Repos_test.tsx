import * as React from "react";
import {Repos} from "sourcegraph/user/settings/Repos";
import testdataData from "sourcegraph/user/settings/testdata/Repos-data.json";
import {autotest} from "sourcegraph/util/testutil/autotest";

describe("Repos", () => {
	it("should render repos", () => {
		let repos = [{
			Private: false,
			URI: "someURL",
			Description: "someDescription",
			UpdatedAt: "2016-02-24T10:18:55-08:00",
			Language: "Go",
		}];
		autotest(testdataData, "sourcegraph/user/settings/testdata/Repos-data.json",
			<Repos repos={repos} />,
		);
	});
});
