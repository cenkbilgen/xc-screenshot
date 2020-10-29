//
//  Screenshot.swift
//  bike-sUITests
//
//  Created by Cenk Bilgen on 2020-10-28.
//  Copyright © 2020 Cenk Bilgen. All rights reserved.
//

import XCTest

extension XCTestCase {
  
  func takeScreenshot(name: String) throws {
    let filename = "screenshot-\(name)-\(Locale.current.description).png"
    let fullScreenshot = XCUIScreen.main.screenshot()
    let screenshot = XCTAttachment(uniformTypeIdentifier: "public.png", name: filename, payload: fullScreenshot.pngRepresentation, userInfo: nil)
    screenshot.lifetime = .keepAlways
    add(screenshot)
    saveLocal(filename: filename, data: fullScreenshot.pngRepresentation)
  }
  
  fileprivate func saveLocal(filename: String, data: Data) {
    var request = URLRequest(url: URL(string: "http://localhost:7000/upload")!)
    request.addValue("filename=\(filename)", forHTTPHeaderField: "Content-Disposition")
    request.httpMethod = "POST"
    let task = URLSession.shared.uploadTask(with: request, from: data) { (_, response, error) in
      print(error.debugDescription)
    }
    task.resume()
  }
  
}


